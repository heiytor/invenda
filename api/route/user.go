package route

import (
	"net/http"

	"github.com/heiytor/invenda/api/pkg/requests"
	"github.com/labstack/echo/v4"
)

func (rs *Routes) userRoutes() []*route {
	return []*route{
		rs.userCreate(),
	}
}

func (rs *Routes) userCreate() *route {
	return &route{
		method: http.MethodPost,
		path:   "/user",
		handler: func(c echo.Context) error {
			req := &requests.UserCreate{}

			if err := c.Bind(req); err != nil {
				return err
			}

			if err := c.Validate(req); err != nil {
				return err
			}

			insertedID, err := rs.service.UserCreate(c.Request().Context(), req)
			if err != nil {
				return err
			}

			c.Response().Header().Set("X-Inserted-Id", insertedID)

			return c.NoContent(http.StatusCreated)
		},
	}
}
