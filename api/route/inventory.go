package route

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (rs *Routes) inventoryRoutes() []*route {
	return []*route{
		rs.inventoryGet(),
		rs.inventoryList(),
		rs.inventoryCreate(),
		rs.inventoryUpdate(),
		rs.inventoryDelete(),
	}
}

func (rs *Routes) inventoryGet() *route {
	return &route{
		method: http.MethodGet,
		path:   "/inventory/:item",
		handler: func(c echo.Context) error {
			return c.NoContent(http.StatusTeapot)
		},
	}
}

func (rs *Routes) inventoryList() *route {
	return &route{
		method: http.MethodGet,
		path:   "/inventory",
		handler: func(c echo.Context) error {
			return c.NoContent(http.StatusTeapot)
		},
	}
}

func (rs *Routes) inventoryCreate() *route {
	return &route{
		method: http.MethodPost,
		path:   "/inventory",
		handler: func(c echo.Context) error {
			return c.NoContent(http.StatusTeapot)
		},
	}
}

func (rs *Routes) inventoryUpdate() *route {
	return &route{
		method: http.MethodPatch,
		path:   "/inventory/:item",
		handler: func(c echo.Context) error {
			return c.NoContent(http.StatusTeapot)
		},
	}
}

func (rs *Routes) inventoryDelete() *route {
	return &route{
		method: http.MethodDelete,
		path:   "/inventory/:item",
		handler: func(c echo.Context) error {
			return c.NoContent(http.StatusTeapot)
		},
	}
}
