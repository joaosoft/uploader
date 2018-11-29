package models

import (
	"github.com/joaosoft/web"
)

type ErrorResponse struct {
	Code    web.Status `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
	Cause   string     `json:"cause,omitempty"`
}

type UploadRequest struct {
	IdUpload string `json:"id_upload"`
	Name     string `json:"name" validate:"nonzero"`
	File     []byte `json:"file" validate:"nonzero"`
}

type UploadResponse struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
