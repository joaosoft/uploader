package storage

import (
	"uploader/models"

	"github.com/joaosoft/errors"
	"github.com/joaosoft/manager"
)

type StorageDatabase struct {
	conn manager.IDB
}

func NewStorageDatabase(connection manager.IDB) *StorageDatabase {
	return &StorageDatabase{
		conn: connection,
	}
}

func (storage *StorageDatabase) Upload(uploadRequest *models.UploadRequest) (*models.UploadResponse, error) {
	_, err := storage.conn.Get().Exec(`
		INSERT INTO uploader.upload(
			id_upload,
			name, 
			file)
		VALUES($1, $2, $3)
	`,
		uploadRequest.IdUpload,
		uploadRequest.Name,
		uploadRequest.File)
	if err != nil {
		return nil, errors.New("upload", err)
	}

	return &models.UploadResponse{
		Name: uploadRequest.Name,
		Path: "",
	}, nil
}
