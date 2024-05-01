package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Binder struct{}

func (b *Binder) Bind(s interface{}, c echo.Context) error {
	binder := new(echo.DefaultBinder)

	if err := binder.Bind(s, c); err != nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}

	return nil
}
