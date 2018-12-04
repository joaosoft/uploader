package storage

import (
	"uploader/models"

	"github.com/joaosoft/logger"
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

	logger.Infof("uploading file %s", uploadRequest.Name)
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
		Name: uploadRequest.Name,
		Path: uploadRequest.IdUpload,
	}, nil
}

func (storage *StorageDatabase) Download(path string) ([]byte, error) {

	logger.Infof("downloading file with path %s", path)
	row := storage.conn.Get().QueryRow(`
		SELECT file FROM upload
		WHERE name = $1
	`,
		path)

	var file []byte
	if err := row.Scan(&file); err != nil {
		return nil, err
	}

	return file, nil
}
