package storage

import (
	"encoding/json"
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
	message, _ := json.Marshal(uploadRequest)
	storage.conn.Publish(routingKey, message, true)

	return &models.UploadResponse{
		IdUpload: uploadRequest.IdUpload,
	}, nil
}

func (storage *StorageRabbitmq) Download(idUpload string) ([]byte, error) {

	logger.Infof("downloading file with id upload %s", idUpload)

	return nil, nil
}
