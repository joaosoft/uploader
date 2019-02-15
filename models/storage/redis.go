package storage

import (
	"fmt"
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type StorageRedis struct {
	conn manager.IRedis
	logger logger.ILogger
}

func NewStorageRedis(connection manager.IRedis, logger logger.ILogger) *StorageRedis {
	return &StorageRedis{
		conn: connection,
		logger: logger,
	}
}

func (storage *StorageRedis) Upload(path string, file []byte) (string, error) {

	storage.logger.Infof("uploading file %s", path)
	key := fmt.Sprintf("image:%s", path)
	storage.conn.Set(key, file)

	return path, nil
}

func (storage *StorageRedis) Download(path string) ([]byte, error) {

	storage.logger.Infof("downloading file with id upload %s", path)
	key := fmt.Sprintf("image:%s", path)

	return storage.conn.Get(key)
}
