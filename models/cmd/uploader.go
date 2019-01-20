package cmd

import (
	"sync"
	"uploader/models"
	"uploader/models/common"
	"uploader/models/controllers"
	"uploader/models/interactors"
	"uploader/models/storage"

	"github.com/labstack/gommon/log"

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
	config, simpleConfig, err := models.NewConfig()
	uploader := &Uploader{
		pm:     manager.NewManager(manager.WithRunInBackground(false)),
		config: &config.Uploader,
	}

	if uploader.isLogExternal {
		uploader.pm.Reconfigure(manager.WithLogger(logger.Instance))
	}

	if err != nil {
		log.Error(err.Error())
	} else {
		uploader.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Uploader.Log.Level)
		logger.Debugf("setting log level to %s", level)
		logger.Reconfigure(logger.WithLevel(level))
	}

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
	simpleDB := manager.NewSimpleDB(&config.Uploader.Db)
	if err := uploader.pm.AddDB("db_postgres", simpleDB); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// redis
	simpleRedis := manager.NewSimpleRedis(&config.Uploader.Redis)
	if err := uploader.pm.AddRedis("redis", simpleRedis); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// rabbitmq
	simpleRabbitmq, err := manager.NewSimpleRabbitmqProducer(&config.Uploader.Rabbitmq)
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
		storageImpl = storage.NewStorageDropbox(dropbox.NewDropbox(dropbox.WithConfiguration(&config.Uploader.Dropbox)))
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
