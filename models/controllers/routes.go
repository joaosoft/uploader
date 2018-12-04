package controllers

import (
	"github.com/joaosoft/manager"
	"github.com/joaosoft/web"
)

func (controller *Controller) RegisterRoutes(w manager.IWeb) error {
	return w.AddRoutes(
		manager.NewRoute(string(web.MethodOptions), "*", controller.DoNothing, web.MiddlewareOptions()),
		manager.NewRoute(string(web.MethodPost), "/api/v1/upload", controller.Upload),
		manager.NewRoute(string(web.MethodGet), "/api/v1/download/:path", controller.Download),
	)
}
