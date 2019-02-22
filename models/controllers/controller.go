package controllers

import (
	"uploader/models"
	"uploader/models/interactors"

	"github.com/joaosoft/errors"
	"github.com/joaosoft/logger"
	"github.com/joaosoft/validator"
	"github.com/joaosoft/web"
)

type Controller struct {
	interactor *interactors.Interactor
	logger     logger.ILogger
}

func NewController(interactor *interactors.Interactor, logger logger.ILogger) *Controller {
	return &Controller{
		interactor: interactor,
		logger:     logger,
	}
}

func (controller *Controller) DoNothing(ctx *web.Context) error {
	return nil
}

func (controller *Controller) Upload(ctx *web.Context) error {
	file := ctx.Request.FormData["file"]
	uploadRequest := &models.UploadRequest{
		Section:  ctx.Request.GetFormDataString("section"),
		Name:     file.Name,
		FileName: file.FileName,
		File:     file.Body,
	}

	if errs := validator.Validate(uploadRequest); len(errs) > 0 {
		err := errors.New(errors.ErrorLevel, 0, "upload", errs)
		controller.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating body of upload request").ToError()
		return ctx.Response.JSON(web.StatusBadRequest, models.ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
	}

	response, err := controller.interactor.Upload(uploadRequest)
	if err != nil {
		err := errors.New(errors.ErrorLevel, 0, err)
		controller.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error uploading file %s", uploadRequest.Name).ToError()
		return ctx.Response.JSON(web.StatusBadRequest, models.ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
	} else {
		return ctx.Response.JSON(web.StatusCreated, response)
	}
}

func (controller *Controller) Download(ctx *web.Context) error {
	downloadRequest := &models.DownloadRequest{
		Section:  ctx.Request.GetUrlParam("section"),
		Size:     ctx.Request.GetUrlParam("size"),
		IdUpload: ctx.Request.GetUrlParam("id_upload"),
	}

	if errs := validator.Validate(downloadRequest); len(errs) > 0 {
		err := errors.New(errors.ErrorLevel, 0, "download", errs)
		controller.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Error("error when validating body of download request").ToError()
		return ctx.Response.JSON(web.StatusBadRequest, models.ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
	}

	response, err := controller.interactor.Download(downloadRequest)
	if err != nil {
		err := errors.New(errors.ErrorLevel, 0, err)
		controller.logger.WithFields(map[string]interface{}{"error": err.Error()}).
			Errorf("error downloading file with id upload %s", downloadRequest.IdUpload).ToError()
		return ctx.Response.JSON(web.StatusBadRequest, models.ErrorResponse{Code: web.StatusBadRequest, Message: err.Error()})
	} else {
		return ctx.Response.File(web.StatusOK, downloadRequest.IdUpload, response)
	}
}
