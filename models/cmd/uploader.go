package cmd

import (
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
	logger        logger.ILogger
	mux           sync.Mutex
}

// NewUploader ...
func NewUploader(options ...UploaderOption) (*Uploader, error) {
	config, simpleConfig, err := models.NewConfig()

	service := &Uploader{
		pm:     manager.NewManager(manager.WithRunInBackground(false)),
		logger: logger.NewLogDefault("uploader", logger.WarnLevel),
		config: &config.Uploader,
	}

	if service.isLogExternal {
		service.pm.Reconfigure(manager.WithLogger(logger.Instance))
	}

	if err != nil {
		service.logger.Error(err.Error())
	} else {
		service.pm.AddConfig("config_app", simpleConfig)
		level, _ := logger.ParseLevel(config.Uploader.Log.Level)
		service.logger.Debugf("setting log level to %s", level)
		service.logger.Reconfigure(logger.WithLevel(level))
	}

	service.Reconfigure(options...)

	// execute migrations
	migration, err := services.NewCmdService(services.WithCmdConfiguration(&service.config.Migration))
	if err != nil {
		return nil, err
	}

	if _, err := migration.Execute(services.OptionUp, 0); err != nil {
		return nil, err
	}

	// database
	simpleDB := manager.NewSimpleDB(&config.Uploader.Db)
	if err := service.pm.AddDB("db_postgres", simpleDB); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// redis
	simpleRedis := manager.NewSimpleRedis(&config.Uploader.Redis)
	if err := service.pm.AddRedis("redis", simpleRedis); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// rabbitmq
	simpleRabbitmq, err := manager.NewSimpleRabbitmqProducer(&config.Uploader.Rabbitmq)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	if err := service.pm.AddRabbitmqProducer("rabbitmq", simpleRabbitmq); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// choose the storage implementation
	var storageImpl interactors.IStorage

	switch service.config.Storage {
	case common.ConstStorageDatabase:
		storageImpl = storage.NewStorageDatabase(simpleDB, service.logger)
	case common.ConstStorageRedis:
		storageImpl = storage.NewStorageRedis(simpleRedis, service.logger)
	case common.ConstStorageRabbitmq:
		storageImpl = storage.NewStorageRabbitmq(simpleRabbitmq, service.logger)
	case common.ConstStorageDropbox:
		storageImpl = storage.NewStorageDropbox(dropbox.NewDropbox(dropbox.WithConfiguration(&config.Uploader.Dropbox)), service.logger)
	}

	// web api
	web := manager.NewSimpleWebServer(service.config.Host)
	controller := controllers.NewController(interactors.NewInteractor(storageImpl, service.logger), service.logger)
	controller.RegisterRoutes(web)

	service.pm.AddWeb("api_web", web)

	return service, nil
}

// Start ...
func (uploader *Uploader) Start() error {
	return uploader.pm.Start()
}

// Stop ...
func (uploader *Uploader) Stop() error {
	return uploader.pm.Stop()
}
