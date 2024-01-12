package route

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var inventoryRoutes = []*route{
	inventoryGet,
	inventoryList,
	inventoryCreate,
	inventoryUpdate,
	inventoryDelete,
}

var inventoryGet = &route{
	method: http.MethodGet,
	path:   "/inventory/:item",
	handler: func(c echo.Context) error {
		return c.NoContent(http.StatusTeapot)
	},
}

var inventoryList = &route{
	method: http.MethodGet,
	path:   "/inventory",
	handler: func(c echo.Context) error {
		return c.NoContent(http.StatusTeapot)
	},
}

var inventoryCreate = &route{
	method: http.MethodPost,
	path:   "/inventory",
	handler: func(c echo.Context) error {
		return c.NoContent(http.StatusTeapot)
	},
}

var inventoryUpdate = &route{
	method: http.MethodPatch,
	path:   "/inventory/:item",
	handler: func(c echo.Context) error {
		return c.NoContent(http.StatusTeapot)
	},
}

var inventoryDelete = &route{
	method: http.MethodDelete,
	path:   "/inventory/:item",
	handler: func(c echo.Context) error {
		return c.NoContent(http.StatusTeapot)
	},
}
