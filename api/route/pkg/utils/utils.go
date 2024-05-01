package utils

import (
	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/labstack/echo/v4"
)

func UserClaims(c echo.Context) *models.UserClaims {
	return c.Get("claims").(*models.UserClaims)
}
