package cmd

import (
	"fmt"
	"sync"
	"uploader/models"
	"uploader/models/common"
	"uploader/models/controllers"
	"uploader/models/interactors"
	"uploader/models/storage"

	"github.com/joaosoft/migration/services"

	"github.com/joaosoft/dropbox"
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type Uploader struct {
	config        *models.UploaderConfig
	isLogExternal bool
	pm            *manager.Manager
	mux           sync.Mutex
}

func init() {
	logger.WithPrefix("service", "uploader")
}

// NewUploader ...
func NewUploader(options ...UploaderOption) (*Uploader, error) {
	uploader := &Uploader{
		pm: manager.NewManager(manager.WithRunInBackground(false)),
	}

	if uploader.isLogExternal {
		uploader.pm.Reconfigure(manager.WithLogger(logger.Instance))
	}

	// load configuration File
	appConfig := &models.AppConfig{}
	if simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", common.GetEnv()), appConfig); err != nil {
		logger.Error(err.Error())
	} else {
		uploader.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(appConfig.Uploader.Log.Level)
		logger.Debugf("setting log level to %s", level)
		logger.Reconfigure(logger.WithLevel(level))
	}

	uploader.config = &appConfig.Uploader

	uploader.Reconfigure(options...)

	if uploader.config.Host == "" {
		uploader.config.Host = common.ConstDefaultURL
	}

	// execute migrations
	migration, err := services.NewCmdService(services.WithCmdConfiguration(&uploader.config.Migration))
	if err != nil {
		return nil, err
	}

	if _, err := migration.Execute(services.OptionUp, 0); err != nil {
		return nil, err
	}

	// database
	simpleDB := manager.NewSimpleDB(&appConfig.Uploader.Db)
	if err := uploader.pm.AddDB("db_postgres", simpleDB); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// redis
	simpleRedis := manager.NewSimpleRedis(&appConfig.Uploader.Redis)
	if err := uploader.pm.AddRedis("redis", simpleRedis); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// rabbitmq
	simpleRabbitmq, err := manager.NewSimpleRabbitmqProducer(&appConfig.Uploader.Rabbitmq)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	if err := uploader.pm.AddRabbitmqProducer("rabbitmq", simpleRabbitmq); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// choose the storage implementation
	var storageImpl interactors.IStorage

	switch uploader.config.Storage {
	case common.ConstStorageDatabase:
		storageImpl = storage.NewStorageDatabase(simpleDB)
	case common.ConstStorageRedis:
		storageImpl = storage.NewStorageRedis(simpleRedis)
	case common.ConstStorageRabbitmq:
		storageImpl = storage.NewStorageRabbitmq(simpleRabbitmq)
	case common.ConstStorageDropbox:
		storageImpl = storage.NewStorageDropbox(dropbox.NewDropbox(dropbox.WithConfiguration(&appConfig.Uploader.Dropbox)))
	}

	// web api
	web := manager.NewSimpleWebServer(uploader.config.Host)
	controller := controllers.NewController(interactors.NewInteractor(storageImpl))
	controller.RegisterRoutes(web)

	uploader.pm.AddWeb("api_web", web)

	return uploader, nil
}

// Start ...
func (uploader *Uploader) Start() error {
	return uploader.pm.Start()
}

// Stop ...
func (uploader *Uploader) Stop() error {
	return uploader.pm.Stop()
}
