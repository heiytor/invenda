package route

import (
	"github.com/heiytor/invenda/api/route/pkg/echoutils"
	"github.com/heiytor/invenda/api/route/pkg/middleware"
	"github.com/heiytor/invenda/api/service"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v3"
)

// route é a definição de uma rota, e é retornado por todos os métodos de rota do Routes.
type route struct {
	method  string           // method é o método HTTP usado pela rota.
	path    string           // path é o caminho a qual a rota deve ser "bindada".
	handler echo.HandlerFunc // handler é o callback a ser executado.
}

// Routes "guarda" todas as rotas da aplicação.
type Routes struct {
	service service.Service // Service é a implementação das regras de negócio.
	E       *echo.Echo
}

func New(service service.Service, logger *lecho.Logger) *Routes {
	e := echo.New()

	e.Logger = logger
	e.Binder = &echoutils.Binder{}
	e.Validator = &echoutils.Validator{}

	e.Use(echomiddleware.RequestID())
	e.Use(middleware.Logger(logger))

	rs := &Routes{E: e, service: service}
	rs.bindRoutes()

	return rs
}

func (rs *Routes) bindRoutes() {
	routes := []*route{}

	routes = append(routes, rs.healthcheck())
	routes = append(routes, rs.inventoryRoutes()...)
	routes = append(routes, rs.userRoutes()...)

	for _, r := range routes {
		log.
			Info().
			Str("method", r.method).
			Str("path", r.path).
			Msg("Registering route")

		rs.E.Add(r.method, r.path, r.handler)
	}
}
