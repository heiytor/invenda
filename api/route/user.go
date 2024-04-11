package route

import (
	"net/http"

	"github.com/heiytor/invenda/api/pkg/requests"
	"github.com/labstack/echo/v4"
)

func (rs *Routes) userRoutes() []*route {
	return []*route{
		rs.userCreate(),
		rs.userUpdate(),
		rs.userAuth(),
	}
}

func (rs *Routes) userCreate() *route {
	return &route{
		method: http.MethodPost,
		path:   "/user",
		handler: func(c echo.Context) error {
			ctx := c.Request().Context()
			req := new(requests.CreateUser)

			if err := c.Bind(req); err != nil {
				return err
			}

			if err := c.Validate(req); err != nil {
				return err
			}

			insertedID, err := rs.service.CreateUser(ctx, req)
			if err != nil {
				return err
			}

			c.Response().Header().Set("X-Inserted-Id", insertedID)
			return c.NoContent(http.StatusCreated)
		},
	}
}

func (rs *Routes) userUpdate() *route {
	return &route{
		method: http.MethodPatch,
		path:   "/user",
		handler: func(c echo.Context) error {
			ctx := c.Request().Context()
			req := new(requests.UpdateUser)

			if err := c.Bind(req); err != nil {
				return err
			}

			if err := c.Validate(req); err != nil {
				return err
			}

			usr, err := rs.service.UpdateUser(ctx, "01HV2ACZ5ZVQ9SC2WE3P8VJ15G", req)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, usr)
		},
	}
}

func (rs *Routes) userAuth() *route {
	return &route{
		method: http.MethodPost,
		path:   "/user/auth",
		handler: func(c echo.Context) error {
			ctx := c.Request().Context()
			req := new(requests.AuthUser)

			if err := c.Bind(req); err != nil {
				return err
			}

			if err := c.Validate(req); err != nil {
				return err
			}

			token, err := rs.service.AuthUser(ctx, req)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, token)
		},
	}
}
