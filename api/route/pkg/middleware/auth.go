package middleware

import (
	"net/http"
	"strings"

	"github.com/heiytor/invenda/api/pkg/errors"
	"github.com/heiytor/invenda/api/pkg/jwt"
	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/labstack/echo/v4"
)

// Auth is responsible for authenticating a request. It expects an "Authorization" header containing a Bearer token.
// If no token is provided, the middleware returns a 401 error. Otherwise, the middleware creates a
// [github.com/heiytor/invenda/api/pkg/models.UserClaims] and stores it in Echo's context.
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization := strings.Split(c.Request().Header.Get("Authorization"), " ")

		if len(authorization) != 2 || authorization[0] != "Bearer" {
			return errors.
				New().
				Code(http.StatusUnauthorized).
				Layer(errors.LayerRoute).
				Msg(errors.MsgInvalidAuthtorization)
		}

		claims := new(models.UserClaims)
		if err := jwt.Decode(authorization[1], claims); err != nil {
			return errors.
				New().
				Code(http.StatusUnauthorized).
				Layer(errors.LayerRoute).
				Msg(errors.MsgInvalidAuthtorization)
		}

		c.Set("claims", claims)

		return next(c)
	}
}
