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
	logger logger.ILogger
}

func NewStorageRabbitmq(connection manager.IRabbitmqProducer, logger logger.ILogger) *StorageRabbitmq {
	return &StorageRabbitmq{
		conn: connection,
		logger: logger,
	}
}

func (storage *StorageRabbitmq) Upload(uploadRequest *models.UploadRequest) (*models.UploadResponse, error) {

	storage.logger.Infof("uploading file %s", uploadRequest.Name)
	routingKey := fmt.Sprintf("upload.file")
	message, _ := json.Marshal(uploadRequest)
	storage.conn.Publish(routingKey, message, true)

	return &models.UploadResponse{
		IdUpload: uploadRequest.IdUpload,
	}, nil
}

func (storage *StorageRabbitmq) Download(idUpload string) ([]byte, error) {

	storage.logger.Infof("downloading file with id upload %s", idUpload)

	return nil, nil
}
