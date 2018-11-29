package interactors

import (
	"github.com/joaosoft/logger"
	"uploader/models"
	"uploader/models/common"
)

type IStorage interface {
	Upload(uploadRequest *models.UploadRequest) (*models.UploadResponse, error)
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
