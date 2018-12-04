package storage

import (
	"fmt"
	"uploader/models"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type StorageRedis struct {
	conn manager.IRedis
}

func NewStorageRedis(connection manager.IRedis) *StorageRedis {
	return &StorageRedis{
		conn: connection,
	}
}

func (storage *StorageRedis) Upload(uploadRequest *models.UploadRequest) (*models.UploadResponse, error) {

	logger.Infof("uploading file %s", uploadRequest.Name)
	key := fmt.Sprintf("image:%s", uploadRequest.Name)
	storage.conn.Set(key, uploadRequest.File)

	return &models.UploadResponse{
		Name: uploadRequest.Name,
		Path: key,
	}, nil
}

func (storage *StorageRedis) Download(path string) ([]byte, error) {

	logger.Infof("downloading file with path %s", path)
	key := fmt.Sprintf("image:%s", path)

	return storage.conn.Get(key)
}
