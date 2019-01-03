package interactors

import (
	"uploader/models"
	"uploader/models/common"

	"github.com/joaosoft/logger"
)

type IStorage interface {
	Upload(uploadRequest *models.UploadRequest) (*models.UploadResponse, error)
	Download(idUpload string) ([]byte, error)
}

type Interactor struct {
	storage IStorage
}

func NewInteractor(storage IStorage) *Interactor {
	return &Interactor{
		storage: storage,
	}
}

func (interactor *Interactor) Upload(uploadRequest *models.UploadRequest) (*models.UploadResponse, error) {
	logger.WithFields(map[string]interface{}{"method": "Upload"})
	logger.Infof("uploading file with name %s", uploadRequest.Name)

	uploadRequest.IdUpload = common.NewULID().String()
	response, err := interactor.storage.Upload(uploadRequest)
	if err != nil {
		logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error uploading file with name %s on storage [ error: %s]", uploadRequest.Name, err).ToError()
	}

	return response, err
}

func (interactor *Interactor) Download(idUpload string) ([]byte, error) {
	logger.WithFields(map[string]interface{}{"method": "Download"})
	logger.Infof("downloading file with id upload %s", idUpload)

	response, err := interactor.storage.Download(idUpload)
	if err != nil {
		logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error uploading file with id upload %s on storage [ error: %s]", idUpload, err).ToError()
	}

	return response, err
}
