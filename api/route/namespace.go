package route

import (
	"net/http"
	"strconv"

	"github.com/heiytor/invenda/api/pkg/auth"
	"github.com/heiytor/invenda/api/pkg/errors"
	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/heiytor/invenda/api/pkg/requests"
	"github.com/heiytor/invenda/api/route/pkg/utils"
	"github.com/labstack/echo/v4"
)

func (rs *Routes) namespaceList() *route[ProtectedHandler] {
	return &route[ProtectedHandler]{
		method:      http.MethodGet,
		path:        "/namespaces",
		group:       GroupPublic,
		middlewares: []echo.MiddlewareFunc{},
		handler: func(c echo.Context, s *models.Session) error {
			ctx := c.Request().Context()
			req := new(requests.ListNamespace)

			if err := c.Bind(req); err != nil {
				return err
			}

			req.Paginator.Normalize()
			req.Sorter.NormalizeWith("created_at")

			if err := c.Validate(req); err != nil {
				return err
			}

			ns, count, err := rs.service.ListNamespace(ctx, s.UserID, req)
			c.Response().Header().Set("X-Total-Count", strconv.FormatInt(count, 10))

			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, ns)
		},
	}
}

func (rs *Routes) namespaceGet() *route[ProtectedHandler] {
	return &route[ProtectedHandler]{
		method:      http.MethodGet,
		path:        "/namespace",
		group:       GroupPublic,
		middlewares: []echo.MiddlewareFunc{},
		handler: func(c echo.Context, s *models.Session) error {
			ctx := c.Request().Context()

			if !auth.Report(s.Permissions, auth.NamespaceRead) {
				return errors.
					New().
					Layer(errors.LayerRoute).
					Attr("required", auth.NamespaceRead).
					Code(http.StatusForbidden).
					Msg(errors.MsgInsufficientPermission)
			}

			ns, err := rs.service.GetNamespace(ctx, s.UserID, s.NamespaceID)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, ns)
		},
	}
}

func (rs *Routes) namespaceCreate() *route[ProtectedHandler] {
	return &route[ProtectedHandler]{
		method:      http.MethodPost,
		path:        "/namespace",
		group:       GroupPublic,
		middlewares: []echo.MiddlewareFunc{},
		handler: func(c echo.Context, s *models.Session) error {
			ctx := c.Request().Context()
			req := new(requests.CreateNamespace)

			if err := c.Bind(req); err != nil {
				return err
			}

			if err := c.Validate(req); err != nil {
				return err
			}

			insertedID, err := rs.service.CreateNamespace(ctx, s.UserID, req)
			if err != nil {
				return err
			}

			c.Response().Header().Set("X-Inserted-Id", insertedID)
			return c.NoContent(http.StatusCreated)
		},
	}
}

func (rs *Routes) namespaceUpdate() *route[ProtectedHandler] {
	return &route[ProtectedHandler]{
		method:      http.MethodPatch,
		path:        "/namespace",
		group:       GroupPublic,
		middlewares: []echo.MiddlewareFunc{},
		handler: func(c echo.Context, s *models.Session) error {
			ctx := c.Request().Context()
			req := new(requests.UpdateNamespace)

			if err := c.Bind(req); err != nil {
				return err
			}

			if err := c.Validate(req); err != nil {
				return err
			}

			if !auth.Report(s.Permissions, auth.NamespaceWrite) {
				return errors.
					New().
					Layer(errors.LayerRoute).
					Attr("required", auth.NamespaceWrite).
					Code(http.StatusForbidden).
					Msg(errors.MsgInsufficientPermission)
			}

			if err := rs.service.UpdateNamespace(ctx, s.UserID, s.NamespaceID, req); err != nil {
				return err
			}

			return c.NoContent(http.StatusOK)
		},
	}
}

func (rs *Routes) namespaceDelete() *route[ProtectedHandler] {
	return &route[ProtectedHandler]{
		method:      http.MethodDelete,
		path:        "/namespace",
		group:       GroupPublic,
		middlewares: []echo.MiddlewareFunc{},
		handler: func(c echo.Context, s *models.Session) error {
			claims := utils.UserClaims(c)
			ctx := c.Request().Context()

			if !auth.Report(claims.Permissions, auth.NamespaceDelete) {
				return errors.
					New().
					Layer(errors.LayerRoute).
					Attr("required", auth.NamespaceDelete).
					Code(http.StatusForbidden).
					Msg(errors.MsgInsufficientPermission)
			}

			if err := rs.service.DeleteNamespace(ctx, s.NamespaceID); err != nil {
				return err
			}

			return c.NoContent(http.StatusNoContent)
		},
	}
}
