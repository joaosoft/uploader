package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/joaosoft/web"
)

type ErrorResponse struct {
	Code    web.Status `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
	Cause   string     `json:"cause,omitempty"`
}

type UploadRequest struct {
	IdUpload string `json:"id_upload" db:"id_upload"`
	Name     string `json:"name" validate:"nonzero" db:"name"`
	Section  string `json:"section" validate:"nonzero" db:"section"`
	FileName string `json:"file_name" validate:"nonzero"`
	File     []byte `json:"file" validate:"nonzero" db:"file"`
}

type UploadResponse struct {
	IdUpload string `json:"id_upload"`
}

type DownloadRequest struct {
	Section  string `json:"section" validate:"nonzero" db:"section"`
	Size     string `json:"size" validate:"nonzero" db:"size"`
	IdUpload string `json:"id_upload" validate:"nonzero"`
}

type Section struct {
	IdSection  int           `json:"id_section" db:"id_section"`
	Name       string        `json:"name" db:"name"`
	Path       string        `json:"path" db:"path"`
	ImageSizes ImageSizeList `json:"image_sizes" db:"image_sizes"`
}

type ImageSize struct {
	Name   string `json:"name" db:"name"`
	Path   string `json:"path" db:"path"`
	Width  int    `json:"width" db:"width"`
	Height int    `json:"height" db:"height"`
}

type ImageSizeList []*ImageSize

func (i *ImageSizeList) Value() (driver.Value, error) {
	j, err := json.Marshal(i)
	return j, err
}
func (i *ImageSizeList) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	fmt.Println(string(src.([]byte)))
	source, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("invalid type")
	}

	err := json.Unmarshal(source, i)
	if err != nil {
		return err
	}

	return nil
}
