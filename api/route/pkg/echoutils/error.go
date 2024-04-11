package echoutils

import (
	"errors"
	"net/http"

	cerrors "github.com/heiytor/invenda/api/pkg/errors"
	"github.com/labstack/echo/v4"
)

// TODO
func ErrorHandler(stdErr error, c echo.Context) {
	var eErr *echo.HTTPError
	var cErr *cerrors.Error

	switch {
	case errors.As(stdErr, &eErr):
		c.NoContent(eErr.Code)
	case errors.As(stdErr, &cErr):
		c.JSON(cErr.Code, cErr)
	default:
		// When an error is not typeof our custom Error or echo.HTTPError, we simply return an
		// Internal Server Error.
		c.NoContent(http.StatusInternalServerError)
	}
}
