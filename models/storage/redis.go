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
	key := fmt.Sprintf("image:%s", uploadRequest.IdUpload)
	storage.conn.Set(key, uploadRequest.File)

	return &models.UploadResponse{
		IdUpload: uploadRequest.IdUpload,
	}, nil
}

func (storage *StorageRedis) Download(idUpload string) ([]byte, error) {

	logger.Infof("downloading file with id upload %s", idUpload)
	key := fmt.Sprintf("image:%s", idUpload)

	return storage.conn.Get(key)
}
