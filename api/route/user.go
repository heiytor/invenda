package route

import (
	"net/http"

	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/heiytor/invenda/api/pkg/requests"
	"github.com/labstack/echo/v4"
)

func (rs *Routes) userRoutes() []*route {
	return []*route{
		rs.userGet(),
		rs.userCreate(),
		rs.userUpdate(),
		rs.userDelete(),
		rs.userAuth(),
	}
}

func (rs *Routes) userGet() *route {
	return &route{
		method:      http.MethodGet,
		path:        "/user/:id",
		protected:   false,
		middlewares: []echo.MiddlewareFunc{},
		handler: func(c echo.Context) error {
			ctx := c.Request().Context()
			req := new(requests.GetUser)

			if err := c.Bind(req); err != nil {
				return err
			}

			if err := c.Validate(req); err != nil {
				return err
			}

			usr, err := rs.service.GetUser(ctx, req)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, usr)
		},
	}
}

func (rs *Routes) userCreate() *route {
	return &route{
		method:      http.MethodPost,
		path:        "/user",
		protected:   false,
		middlewares: []echo.MiddlewareFunc{},
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
		method:      http.MethodPatch,
		path:        "/user",
		protected:   true,
		middlewares: []echo.MiddlewareFunc{},
		handler: func(c echo.Context) error {
			ctx := c.Request().Context()
			req := new(requests.UpdateUser)

			if err := c.Bind(req); err != nil {
				return err
			}

			if err := c.Validate(req); err != nil {
				return err
			}

			claims := c.Get("claims").(*models.UserClaims)
			usr, err := rs.service.UpdateUser(ctx, claims.Subject, req)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, usr)
		},
	}
}

func (rs *Routes) userDelete() *route {
	return &route{
		method:      http.MethodDelete,
		path:        "/user",
		protected:   true,
		middlewares: []echo.MiddlewareFunc{},
		handler: func(c echo.Context) error {
			ctx := c.Request().Context()

			claims := c.Get("claims").(*models.UserClaims)
			if err := rs.service.DeleteUser(ctx, claims.Subject); err != nil {
				return err
			}

			return c.NoContent(http.StatusNoContent)
		},
	}
}

func (rs *Routes) userAuth() *route {
	return &route{
		method:      http.MethodPost,
		path:        "/user/auth",
		middlewares: []echo.MiddlewareFunc{},
		handler: func(c echo.Context) error {
			ctx := c.Request().Context()
			req := new(requests.AuthUser)

			if err := c.Bind(req); err != nil {
				return err
			}

			if err := c.Validate(req); err != nil {
				return err
			}

			claims, token, err := rs.service.AuthUser(ctx, req)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, map[string]interface{}{"token": token, "claims": claims})
		},
	}
}
