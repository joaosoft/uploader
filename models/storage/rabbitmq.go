package storage

import (
	"fmt"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type StorageRabbitmq struct {
	conn   manager.IRabbitmqProducer
	logger logger.ILogger
}

func NewStorageRabbitmq(connection manager.IRabbitmqProducer, logger logger.ILogger) *StorageRabbitmq {
	return &StorageRabbitmq{
		conn:   connection,
		logger: logger,
	}
}

func (storage *StorageRabbitmq) Upload(path string, file []byte) (string, error) {

	storage.logger.Infof("uploading file %s", path)
	routingKey := fmt.Sprintf("upload.file")
	storage.conn.Publish(routingKey, file, true)

	return path, nil
}

func (storage *StorageRabbitmq) Download(path string) ([]byte, error) {

	storage.logger.Infof("downloading file with id upload %s", path)

	return nil, nil
}
