package controllers

import (
	"encoding/base64"
	"uploader/models"
	"uploader/models/interactors"

	"github.com/joaosoft/errors"
	"github.com/joaosoft/logger"
	"github.com/joaosoft/validator"
	"github.com/joaosoft/web"
)

type Controller struct {
	interactor *interactors.Interactor
}

func NewController(interactor *interactors.Interactor) *Controller {
	return &Controller{
		interactor: interactor,
	}
}

func (controller *Controller) DoNothing(ctx *web.Context) error {
	return nil
}

func (controller *Controller) Upload(ctx *web.Context) error {
	file := ctx.Request.Attachments["file"]
	uploadRequest := &models.UploadRequest{
		Name: file.File,
		File: make([]byte, base64.StdEncoding.EncodedLen(len(file.Body))),
	}

	base64.StdEncoding.Encode(uploadRequest.File, file.Body)

	if errs := validator.Validate(uploadRequest); len(errs) > 0 {
		err := errors.New("upload", errs)
		logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating body uploadRequest").ToError()
		return ctx.Response.JSON(web.StatusBadRequest, models.ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
	}

	response, err := controller.interactor.Upload(uploadRequest)
	if err != nil {
		err := errors.New("0", err)
		logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error uploading file %s", uploadRequest.Name).ToError()
		return ctx.Response.JSON(web.StatusBadRequest, models.ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
	} else {
		return ctx.Response.JSON(web.StatusCreated, response)
	}
}

func (controller *Controller) Download(ctx *web.Context) error {
	path := ctx.Request.GetUrlParam("path")
	downloadRequest := &models.DownloadRequest{
		Path: path,
	}

	if errs := validator.Validate(downloadRequest); len(errs) > 0 {
		err := errors.New("download", errs)
		logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating body downloadRequest").ToError()
		return ctx.Response.JSON(web.StatusBadRequest, models.ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
	}

	response, err := controller.interactor.Download(path)
	if err != nil {
		err := errors.New("0", err)
		logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error downloading file with path %s", path).ToError()
		return ctx.Response.JSON(web.StatusBadRequest, models.ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
	} else {
		decoded := make([]byte, base64.StdEncoding.DecodedLen(len(response)))
		base64.StdEncoding.Decode(decoded, response)

		return ctx.Response.File(web.StatusOK, path, decoded)
	}
}
