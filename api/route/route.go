package route

import (
	"github.com/heiytor/invenda/api/route/pkg/echoutils"
	"github.com/heiytor/invenda/api/route/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/ziflex/lecho/v3"
)

type route struct {
	method  string
	path    string
	handler echo.HandlerFunc
}

func New(logger *lecho.Logger) *echo.Echo {
	e := echo.New()

	e.Logger = logger
	e.Binder = &echoutils.Binder{}
	e.Validator = &echoutils.Validator{}

	e.Use(middleware.Logger(logger))

	retisterPublicRoutes(e)

	return e
}

func retisterPublicRoutes(e *echo.Echo) {
	public := e.Group("/api")
	publicRoutes := []*route{}

	publicRoutes = append(publicRoutes, healthCheck)
	publicRoutes = append(publicRoutes, inventoryRoutes...)

	for _, r := range publicRoutes {
		public.Add(r.method, r.path, r.handler)
	}
}
