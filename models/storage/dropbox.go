package storage

import (
	"fmt"
	"uploader/models"

	"github.com/joaosoft/dropbox"
	"github.com/joaosoft/logger"
)

type StorageDropbox struct {
	conn *dropbox.Dropbox
}

func NewStorageDropbox(connection *dropbox.Dropbox) *StorageDropbox {
	return &StorageDropbox{
		conn: connection,
	}
}

func (storage *StorageDropbox) Upload(uploadRequest *models.UploadRequest) (*models.UploadResponse, error) {

	logger.Infof("uploading file %s", uploadRequest.Name)
	response, err := storage.conn.File().Upload(fmt.Sprintf("/%s", uploadRequest.Name), uploadRequest.File)
	if err != nil {
		return nil, err
	}

	return &models.UploadResponse{
		Name: uploadRequest.Name,
		Path: response.PathDisplay,
	}, nil
}

func (storage *StorageDropbox) Download(path string) ([]byte, error) {

	logger.Infof("downloading file with path %s", path)
	response, err := storage.conn.File().Download(fmt.Sprintf("/%s", path))
	if err != nil {
		return nil, err
	}

	return response, nil
}
