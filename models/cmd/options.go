package cmd

import (
	"uploader/models"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

// UploaderOption ...
type UploaderOption func(client *Uploader)

// Reconfigure ...
func (uploader *Uploader) Reconfigure(options ...UploaderOption) {
	for _, option := range options {
		option(uploader)
	}
}

// WithConfiguration ...
func WithConfiguration(config *models.UploaderConfig) UploaderOption {
	return func(client *Uploader) {
		client.config = config
	}
}

// WithLogger ...
func WithLogger(l logger.ILogger) UploaderOption {
	return func(uploader *Uploader) {
		logger.Instance = l
		uploader.isLogExternal = true
	}
}

// WithLogLevel ...
func WithLogLevel(level logger.Level) UploaderOption {
	return func(uploader *Uploader) {
		logger.SetLevel(level)
	}
}

// WithManager ...
func WithManager(mgr *manager.Manager) UploaderOption {
	return func(uploader *Uploader) {
		uploader.pm = mgr
	}
}
