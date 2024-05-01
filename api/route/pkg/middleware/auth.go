package middleware

import (
	"strings"

	"github.com/heiytor/invenda/api/pkg/errors"
	"github.com/heiytor/invenda/api/pkg/jwt"
	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/labstack/echo/v4"
)

// TODO: maybe we can remove this middleware from here and put in a gateway (like nginx)?
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization := strings.Split(c.Request().Header.Get("Authorization"), " ")

		if len(authorization) != 2 || authorization[0] != "Bearer" {
			return errors.
				New().
				Code(401).
				Layer(errors.LayerRoute).
				Msg(errors.MsgInvalidAuthtorization)
		}

		claims := new(models.UserClaims)
		if err := jwt.Decode(authorization[1], claims); err != nil {
			return errors.
				New().
				Code(401).
				Layer(errors.LayerRoute).
				Msg(errors.MsgInvalidAuthtorization)
		}

		c.Set("claims", claims)

		return next(c)
	}
}
