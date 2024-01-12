package route

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var healthCheck = &route{
	method: http.MethodGet,
	path:   "/healthcheck",
	handler: func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	},
}
