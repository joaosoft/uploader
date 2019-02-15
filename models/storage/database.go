package storage

import (
	"uploader/models"

	"github.com/joaosoft/dbr"
	"github.com/joaosoft/logger"
	"github.com/joaosoft/manager"
)

type StorageDatabase struct {
	conn   manager.IDB
	logger logger.ILogger
	db     *dbr.Dbr
}

func NewStorageDatabase(connection manager.IDB, dialect string, logger logger.ILogger) (*StorageDatabase, error) {
	db, err := dbr.New(dbr.WithDatabase(dbr.NewDb(connection.Get(), dbr.NewDialect(dialect))))
	if err != nil {
		return nil, err
	}
	return &StorageDatabase{
		conn:   connection,
		logger: logger,
		db:     db,
	}, nil
}

func (storage *StorageDatabase) Upload(path string, file []byte) (string, error) {

	storage.logger.Infof("uploading file %s", path)
	_, err := storage.db.Insert().
		Into("uploader.upload").
		Columns([]interface{}{"id_upload", "file"}...).
		Values(path, file).Exec()

	if err != nil {
		return "", err
	}

	return path, nil
}

func (storage *StorageDatabase) Download(path string) ([]byte, error) {

	storage.logger.Infof("downloading file with id upload %s", path)

	var file []byte
	_, err := storage.db.Select([]interface{}{"file"}...).
		From("uploader.upload").
		Where("id_upload = ?", path).
		Load(&file)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (storage *StorageDatabase) LoadSections() ([]*models.Section, error) {

	storage.logger.Infof("loading sections")

	sections := make([]*models.Section, 0)

	_, err := storage.db.Select([]interface{}{
		"s.id_section",
		"s.name",
		"s.path",
		`COALESCE((SELECT ARRAY_TO_JSON(ARRAY_AGG(ROW_TO_JSON(t))) 
				FROM (SELECT "name",
			   		   "path",
			   		   width
			   		   hight
			  	FROM uploader.section_image_size sisize
			  	JOIN uploader.image_size isize ON isize.id_image_size = sisize.fk_image_size
			  	WHERE sisize.fk_section = s.id_section) t), '[]') AS image_sizes`}...).
		From(dbr.As("uploader.section", "s")).
		Where("s.active").
		Load(&sections)

	if err != nil {
		return nil, err
	}

	return sections, nil
}
