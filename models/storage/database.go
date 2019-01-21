package storage

import (
	"uploader/models"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type StorageDatabase struct {
	conn   manager.IDB
	logger logger.ILogger
}

func NewStorageDatabase(connection manager.IDB, logger logger.ILogger) *StorageDatabase {
	return &StorageDatabase{
		conn:   connection,
		logger: logger,
	}
}

func (storage *StorageDatabase) Upload(uploadRequest *models.UploadRequest) (*models.UploadResponse, error) {

	storage.logger.Infof("uploading file %s", uploadRequest.Name)
	_, err := storage.conn.Get().Exec(`
		INSERT INTO upload(
			id_upload,
			name, 
			file)
		VALUES($1, $2, $3)
	`,
		uploadRequest.IdUpload,
		uploadRequest.Name,
		uploadRequest.File)
	if err != nil {
		return nil, err
	}

	return &models.UploadResponse{
		IdUpload: uploadRequest.IdUpload,
	}, nil
}

func (storage *StorageDatabase) Download(idUpload string) ([]byte, error) {

	storage.logger.Infof("downloading file with id upload %s", idUpload)
	row := storage.conn.Get().QueryRow(`
		SELECT file FROM upload
		WHERE id_upload = $1
	`,
		idUpload)

	var file []byte
	if err := row.Scan(&file); err != nil {
		return nil, err
	}

	return file, nil
}
