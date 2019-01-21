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
	logger  logger.ILogger
}

func NewInteractor(storage IStorage, logger logger.ILogger) *Interactor {
	return &Interactor{
		storage: storage,
		logger:  logger,
	}
}

func (interactor *Interactor) Upload(uploadRequest *models.UploadRequest) (*models.UploadResponse, error) {
	interactor.logger.WithFields(map[string]interface{}{"method": "Upload"})
	interactor.logger.Infof("uploading file with name %s", uploadRequest.Name)

	uploadRequest.IdUpload = common.NewULID().String()
	response, err := interactor.storage.Upload(uploadRequest)
	if err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error uploading file with name %s on storage [ error: %s]", uploadRequest.Name, err).ToError()
	}

	return response, err
}

func (interactor *Interactor) Download(idUpload string) ([]byte, error) {
	interactor.logger.WithFields(map[string]interface{}{"method": "Download"})
	interactor.logger.Infof("downloading file with id upload %s", idUpload)

	response, err := interactor.storage.Download(idUpload)
	if err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error uploading file with id upload %s on storage [ error: %s]", idUpload, err).ToError()
	}

	return response, err
}
