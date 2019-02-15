package interactors

import (
	"fmt"
	"uploader/models"
	"uploader/models/common"
	"uploader/models/storage"

	"github.com/joaosoft/logger"
	"github.com/joaosoft/pictures"
)

type IStorage interface {
	Upload(path string, file []byte) (string, error)
	Download(path string) ([]byte, error)
}

type Interactor struct {
	storage  IStorage
	database *storage.StorageDatabase
	logger   logger.ILogger
	sections map[string]*models.Section
}

func NewInteractor(database *storage.StorageDatabase, storage IStorage, logger logger.ILogger) (*Interactor, error) {

	// load sections
	sectionsList, err := database.LoadSections()
	if err != nil {
		return nil, err
	}

	sectionsMap := make(map[string]*models.Section)
	for _, section := range sectionsList {
		sectionsMap[section.Name] = section
	}

	return &Interactor{
		storage:  storage,
		logger:   logger,
		database: database,
		sections: sectionsMap,
	}, nil
}

func (interactor *Interactor) Upload(uploadRequest *models.UploadRequest) (*models.UploadResponse, error) {
	interactor.logger.WithFields(map[string]interface{}{"method": "Upload"})
	interactor.logger.Infof("uploading file with name %s", uploadRequest.Name)

	// upload the original
	uploadRequest.IdUpload = common.NewULID().String()
	_, err := interactor.storage.Upload(fmt.Sprintf("%s/%s", uploadRequest.Section, "original"), uploadRequest.File)
	if err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error uploading original file with name %s on storage [ error: %s]", uploadRequest.Name, err).ToError()
	}

	go func() {
		for _, section := range interactor.sections {
			for _, imageSize := range section.ImageSizes {
				picture, err := pictures.FromBytes(uploadRequest.File)
				if err != nil {
					interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
						Errorf("error uploading file with size (width %d, height %d) with name %s on storage [ error: %s]", imageSize.Width, imageSize.Height, uploadRequest.Name, err).ToError()
				}

				picture.Resize(imageSize.Width, imageSize.Height)

				uploadRequest.File, err = picture.ToBytes()
				if err != nil {
					interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
						Errorf("error uploading file with size (width %d, height %d) with name %s on storage [ error: %s]", imageSize.Width, imageSize.Height, uploadRequest.Name, err).ToError()
				}

				_, err = interactor.storage.Upload(fmt.Sprintf("%s/%s", section.Path, imageSize.Path), uploadRequest.File)
				if err != nil {
					interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
						Errorf("error uploading file with size (width %d, height %d) with name %s on storage [ error: %s]", imageSize.Width, imageSize.Height, uploadRequest.Name, err).ToError()
				}
			}
		}
	}()

	return &models.UploadResponse{IdUpload: uploadRequest.IdUpload}, err
}

func (interactor *Interactor) Download(downloadRequest *models.DownloadRequest) ([]byte, error) {
	interactor.logger.WithFields(map[string]interface{}{"method": "Download"})
	interactor.logger.Infof("downloading file with id upload %s", downloadRequest.IdUpload)

	response, err := interactor.storage.Download(fmt.Sprintf("%s/%s/%s", downloadRequest.Section, downloadRequest.Size, downloadRequest.IdUpload))
	if err != nil {
		interactor.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error uploading file with id upload %s on storage [ error: %s]", downloadRequest.IdUpload, err).ToError()
	}

	return response, err
}
