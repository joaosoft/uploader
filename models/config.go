package models

import (
	"fmt"
	"uploader/models/common"
	"github.com/joaosoft/manager"
	"github.com/joaosoft/logger"
)

// AppConfig ...
type AppConfig struct {
	Uploader UploaderConfig `json:"uploader"`
}

// UploaderConfig ...
type UploaderConfig struct {
	Host string           `json:"host"`
	Db   manager.DBConfig `json:"db"`
	Log  struct {
		Level string `json:"level"`
	} `json:"log"`
}

// NewConfig ...
func NewConfig(host string, db manager.DBConfig) *UploaderConfig {
	appConfig := &AppConfig{}
	if _, err := manager.NewSimpleConfig(fmt.Sprintf("/config/app.%s.json", common.GetEnv()), appConfig); err != nil {
		logger.Error(err.Error())
	}

	appConfig.Uploader.Host = host

	return &appConfig.Uploader
}
