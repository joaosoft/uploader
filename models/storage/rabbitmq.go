package storage

import (
	"fmt"
	"uploader/models"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type StorageRabbitmq struct {
	conn manager.IRabbitmqProducer
}

func NewStorageRabbitmq(connection manager.IRabbitmqProducer) *StorageRabbitmq {
	return &StorageRabbitmq{
		conn: connection,
	}
}

func (storage *StorageRabbitmq) Upload(uploadRequest *models.UploadRequest) (*models.UploadResponse, error) {

	logger.Infof("uploading file %s", uploadRequest.Name)
	routingKey := fmt.Sprintf("upload.file")
	storage.conn.Publish(routingKey, uploadRequest.File, true)

	return &models.UploadResponse{
		Name: uploadRequest.Name,
		Path: routingKey,
	}, nil
}

func (storage *StorageRabbitmq) Download(path string) ([]byte, error) {

	logger.Infof("downloading file with path %s", path)

	return nil, nil
}
