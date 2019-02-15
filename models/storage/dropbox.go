package storage

import (
	"fmt"

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

func (storage *StorageDropbox) Upload(path string, file []byte) (string, error) {

	storage.logger.Infof("uploading file %s", path)
	response, err := storage.conn.File().Upload(fmt.Sprintf("/%s", path), file)
	if err != nil {
		return "", err
	}

	return response.PathDisplay, nil
}

func (storage *StorageDropbox) Download(path string) ([]byte, error) {

	storage.logger.Infof("downloading file with id upload %s", path)
	response, err := storage.conn.File().Download(fmt.Sprintf("/%s", path))
	if err != nil {
		return nil, err
	}

	return response, nil
}
