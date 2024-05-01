package route

import (
	"github.com/heiytor/invenda/api/pkg/validator"
	"github.com/heiytor/invenda/api/route/pkg/echoutils"
	"github.com/heiytor/invenda/api/service"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

// route é a definição de uma rota, e é retornado por todos os métodos de rota do Routes.
type route struct {
	method      string // method é o método HTTP usado pela rota.
	path        string // path é o caminho a qual a rota deve ser "bindada".
	middlewares []echo.MiddlewareFunc
	handler     echo.HandlerFunc // handler é o callback a ser executado.
}

// Routes "guarda" todas as rotas da aplicação.
type Routes struct {
	service service.Service // Service é a implementação das regras de negócio.
	E       *echo.Echo
}

func New(service service.Service) *Routes {
	r := &Routes{E: echo.New(), service: service}

	r.E.Binder = &echoutils.Binder{}
	r.E.Validator = validator.New()
	r.E.HTTPErrorHandler = echoutils.ErrorHandler

	r.E.Use(echomiddleware.RequestID())

	r.bindRoutes()

	return r
}

func (rs *Routes) bindRoutes() {
	// healthcheck is used by docker to check if the container is healthy
	rs.E.GET("/healthcheck", func(c echo.Context) error {
		return c.NoContent(200)
	})

	routes := []*route{}
	routes = append(routes, rs.userRoutes()...)

	for _, r := range routes {
		log.Info().
			Str("method", r.method).
			Str("path", r.path).
			Msg("Registering route")

		rs.E.Add(r.method, r.path, r.handler, r.middlewares...)
	}
}
