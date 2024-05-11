package route

import (
	"net/http"

	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/heiytor/invenda/api/pkg/requests"
	"github.com/labstack/echo/v4"
)

func (rs *Routes) userGet() *route[echo.HandlerFunc] {
	return &route[echo.HandlerFunc]{
		method:      http.MethodGet,
		path:        "/user/:id",
		group:       GroupPublic,
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

func (rs *Routes) userCreate() *route[echo.HandlerFunc] {
	return &route[echo.HandlerFunc]{
		method:      http.MethodPost,
		path:        "/user",
		group:       GroupPublic,
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

func (rs *Routes) userUpdate() *route[ProtectedHandler] {
	return &route[ProtectedHandler]{
		method:      http.MethodPatch,
		path:        "/user",
		group:       GroupPublic,
		middlewares: []echo.MiddlewareFunc{},
		handler: func(c echo.Context, s *models.Session) error {
			ctx := c.Request().Context()
			req := new(requests.UpdateUser)

			if err := c.Bind(req); err != nil {
				return err
			}

			if err := c.Validate(req); err != nil {
				return err
			}

			usr, err := rs.service.UpdateUser(ctx, s.UserID, req)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, usr)
		},
	}
}

func (rs *Routes) userDelete() *route[ProtectedHandler] {
	return &route[ProtectedHandler]{
		method:      http.MethodDelete,
		path:        "/user",
		group:       GroupPublic,
		middlewares: []echo.MiddlewareFunc{},
		handler: func(c echo.Context, s *models.Session) error {
			ctx := c.Request().Context()

			if err := rs.service.DeleteUser(ctx, s.UserID); err != nil {
				return err
			}

			return c.NoContent(http.StatusNoContent)
		},
	}
}

func (rs *Routes) userCreateSession() *route[echo.HandlerFunc] {
	return &route[echo.HandlerFunc]{
		method:      http.MethodPost,
		path:        "/user/session",
		group:       GroupPublic,
		middlewares: []echo.MiddlewareFunc{},
		handler: func(c echo.Context) error {
			ctx := c.Request().Context()
			req := new(requests.CreateSession)

			if err := c.Bind(req); err != nil {
				return err
			}

			if err := c.Validate(req); err != nil {
				return err
			}

			insertedID, err := rs.service.CreateSession(ctx, req)
			if err != nil {
				return err
			}

			c.Response().Header().Set("X-Inserted-Id", insertedID)
			return c.NoContent(http.StatusOK)
		},
	}
}

func (rs *Routes) userUpdateSession() *route[echo.HandlerFunc] {
	return &route[echo.HandlerFunc]{
		method:      http.MethodPut,
		path:        "/user/session/:id",
		group:       GroupPublic,
		middlewares: []echo.MiddlewareFunc{},
		handler: func(c echo.Context) error {
			ctx := c.Request().Context()
			req := new(requests.UpdateSession)

			if err := c.Bind(req); err != nil {
				return err
			}

			if err := c.Validate(req); err != nil {
				return err
			}

			if err := rs.service.UpdateSession(ctx, req); err != nil {
				return err
			}

			return c.NoContent(http.StatusOK)
		},
	}
}
