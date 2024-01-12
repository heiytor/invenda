package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/ziflex/lecho/v3"
)

func Logger(logger *lecho.Logger) echo.MiddlewareFunc {
	return lecho.Middleware(lecho.Config{Logger: logger})
}
