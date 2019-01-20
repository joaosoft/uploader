package models

import (
	"fmt"
	"uploader/models/common"

	"github.com/joaosoft/migration/services"

	"github.com/joaosoft/dropbox"

	"github.com/joaosoft/manager"
)

// AppConfig ...
type AppConfig struct {
	Uploader UploaderConfig `json:"uploader"`
}

// UploaderConfig ...
type UploaderConfig struct {
	Storage   string                   `json:"storage"`
	Host      string                   `json:"host"`
	Db        manager.DBConfig         `json:"db"`
	Redis     manager.RedisConfig      `json:"redis"`
	Rabbitmq  manager.RabbitmqConfig   `json:"rabbitmq"`
	Dropbox   dropbox.DropboxConfig    `json:"dropbox"`
	Migration services.MigrationConfig `json:"migration"`
	Log       struct {
		Level string `json:"level"`
	} `json:"log"`
}

// NewConfig ...
func NewConfig() (*AppConfig, manager.IConfig, error) {
	appConfig := &AppConfig{}
	simpleConfig, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", common.GetEnv()), appConfig)

	return appConfig, simpleConfig, err
}
