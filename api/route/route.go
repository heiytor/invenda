package route

import (
	"net/http"

	"github.com/heiytor/invenda/api/pkg/cache"
	"github.com/heiytor/invenda/api/pkg/models"
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

type Handler echo.HandlerFunc
type ProtectedHandler func(c echo.Context, s *models.Session) error

type route[T any] struct {
	method      string
	path        string
	group       Group
	middlewares []echo.MiddlewareFunc
	handler     T
}

type Routes struct {
	service service.Service
	E       *echo.Echo
}

func New(service service.Service, cache cache.Cache) *Routes {
	r := &Routes{E: echo.New(), service: service}

	r.E.Binder = &utils.Binder{}
	r.E.Validator = validator.New()
	r.E.HTTPErrorHandler = utils.ErrorHandler

	r.E.Use(echomiddleware.RequestID())

	// healthcheck is used by docker to check if the container is healthy
	r.E.GET("/healthcheck", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	handlers, protectedHandlers := r.allRoutes()
	pub := r.E.Group(string(GroupPublic))
	pri := r.E.Group(string(GroupInternal))

	for _, h := range handlers {
		log.Info().
			Str("method", h.method).
			Str("path", h.path).
			Str("group", string(h.group)).
			Msg("Registering non-protected route")

		switch h.group {
		case GroupPublic:
			pub.Add(h.method, h.path, h.handler, h.middlewares...)
		case GroupInternal:
			pri.Add(h.method, h.path, h.handler, h.middlewares...)
		}
	}

	for _, h := range protectedHandlers {
		log.Info().
			Str("method", h.method).
			Str("path", h.path).
			Str("group", string(h.group)).
			Msg("Registering protected route")

		switch h.group {
		case GroupPublic:
			pub.Add(h.method, h.path, middleware.Auth(cache, h.handler), h.middlewares...)
		case GroupInternal:
			pri.Add(h.method, h.path, middleware.Auth(cache, h.handler), h.middlewares...)
		}
	}

	return r
}

func (rs *Routes) allRoutes() ([]*route[echo.HandlerFunc], []*route[ProtectedHandler]) {
	handlers := []*route[echo.HandlerFunc]{
		rs.userGet(),
		rs.userCreate(),
		rs.userCreateSession(),
		rs.userUpdateSession(),
	}

	protectedHandlers := []*route[ProtectedHandler]{
		rs.userUpdate(),
		rs.userDelete(),

		rs.namespaceGet(),
		rs.namespaceList(),
		rs.namespaceCreate(),
		rs.namespaceUpdate(),
		rs.namespaceDelete(),
	}

	return handlers, protectedHandlers
}
