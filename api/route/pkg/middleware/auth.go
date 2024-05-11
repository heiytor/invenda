package middleware

import (
	"net/http"
	"strings"

	"github.com/heiytor/invenda/api/pkg/auth"
	"github.com/heiytor/invenda/api/pkg/cache"
	"github.com/heiytor/invenda/api/pkg/errors"
	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/labstack/echo/v4"
)

// Auth is responsible for authenticating a request. It expects an "Authorization" header containing a Bearer token.
// If no token is provided, the middleware returns a 401 error. Otherwise, the middleware creates a
// [github.com/heiytor/invenda/api/pkg/models.UserClaims] and stores it in Echo's context.
func Auth(cache cache.Cache, handler func(c echo.Context, s *models.Session) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Request().Header.Get("X-Session-ID")
		if id == "" {
			return errors.
				New().
				Code(http.StatusUnauthorized).
				Layer(errors.LayerRoute).
				Msg(errors.MsgInvalidAuthtorization)
		}

		str := ""
		if err := cache.Get(c.Request().Context(), id, &str); err != nil || str == "" {
			return errors.
				New().
				Code(http.StatusUnauthorized).
				Layer(errors.LayerRoute).
				Msg(errors.MsgInvalidAuthtorization)
		}

		parts := strings.Split(str, ";")
		session := &models.Session{
			NamespaceID: parts[0],
			UserID:      parts[1],
			Permissions: auth.Permissions{}.FromString(parts[2]),
		}

		return handler(c, session)
	}
}
