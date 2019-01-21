package storage

import (
	"fmt"
	"uploader/models"

	"github.com/joaosoft/dropbox"
	"github.com/joaosoft/logger"
)

type StorageDropbox struct {
	conn   *dropbox.Dropbox
	logger logger.ILogger
}

func NewStorageDropbox(connection *dropbox.Dropbox, logger logger.ILogger) *StorageDropbox {
	return &StorageDropbox{
		conn:   connection,
		logger: logger,
	}
}

func (storage *StorageDropbox) Upload(uploadRequest *models.UploadRequest) (*models.UploadResponse, error) {

	storage.logger.Infof("uploading file %s", uploadRequest.Name)
	response, err := storage.conn.File().Upload(fmt.Sprintf("/%s", uploadRequest.IdUpload), uploadRequest.File)
	if err != nil {
		return nil, err
	}

	return &models.UploadResponse{
		IdUpload: response.PathDisplay,
	}, nil
}

func (storage *StorageDropbox) Download(idUpload string) ([]byte, error) {

	storage.logger.Infof("downloading file with id upload %s", idUpload)
	response, err := storage.conn.File().Download(fmt.Sprintf("/%s", idUpload))
	if err != nil {
		return nil, err
	}

	return response, nil
}
