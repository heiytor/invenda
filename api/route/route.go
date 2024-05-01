package route

import (
	"net/http"

	"github.com/heiytor/invenda/api/pkg/validator"
	"github.com/heiytor/invenda/api/route/pkg/middleware"
	"github.com/heiytor/invenda/api/route/pkg/utils"
	"github.com/heiytor/invenda/api/service"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

type Group string

const (
	GroupPublic   Group = "public"
	GroupInternal Group = "internal"
)

type route struct {
	method      string
	path        string
	protected   bool
	group       Group
	middlewares []echo.MiddlewareFunc
	handler     echo.HandlerFunc
}

type Routes struct {
	service service.Service
	E       *echo.Echo
}

func New(service service.Service) *Routes {
	r := &Routes{E: echo.New(), service: service}

	r.E.Binder = &utils.Binder{}
	r.E.Validator = validator.New()
	r.E.HTTPErrorHandler = utils.ErrorHandler

	r.E.Use(echomiddleware.RequestID())

	r.bindRoutes()

	return r
}

func (rs *Routes) bindRoutes() {
	// healthcheck is used by docker to check if the container is healthy
	rs.E.GET("/healthcheck", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	routes := []*route{}
	routes = append(routes, rs.userRoutes()...)
	routes = append(routes, rs.namespaceRoutes()...)

	pub := rs.E.Group(string(GroupPublic))
	pri := rs.E.Group(string(GroupInternal))

	for _, r := range routes {
		log.Info().
			Str("method", r.method).
			Str("path", r.path).
			Bool("protected", r.protected).
			Str("group", string(r.group)).
			Msg("Registering route")

		if r.protected {
			r.middlewares = append([]echo.MiddlewareFunc{middleware.Auth}, r.middlewares...)
		}

		switch r.group {
		case GroupPublic:
			pub.Add(r.method, r.path, r.handler, r.middlewares...)
		case GroupInternal:
			pri.Add(r.method, r.path, r.handler, r.middlewares...)
		}
	}
}
