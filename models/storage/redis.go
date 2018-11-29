package storage

import (
	"uploader/models"

	"github.com/joaosoft/manager"
)

type StorageRedis struct {
	conn manager.IDB
}

func NewStorageRedis(connection manager.IDB) *StorageRedis {
	return &StorageRedis{
		conn: connection,
	}
}

func (storage *StorageRedis) Upload(uploadRequest *models.UploadRequest) (*models.UploadResponse, error) {

	return &models.UploadResponse{
		Name: uploadRequest.Name,
		Path: "",
	}, nil
}
