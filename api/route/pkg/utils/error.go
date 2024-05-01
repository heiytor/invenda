package utils

import (
	"errors"
	"net/http"

	cerrors "github.com/heiytor/invenda/api/pkg/errors"
	"github.com/labstack/echo/v4"
)

// ErrorHandler manages errors returned from the routing layer.
// Echo errors, such as "route not found," return a response with no content body and the specified status code.
// Custom errors from the pkg/errors package return a response with the specified body and status code.
// All other errors, such as MongoDB errors and unhandled errors, return a response with no content body and a status code of 500.
// These errors are also reported to Sentry.
func ErrorHandler(err error, c echo.Context) {
	var eErr *echo.HTTPError
	var cErr *cerrors.Error

	switch {
	case errors.As(err, &eErr):
		c.NoContent(eErr.Code)
	case errors.As(err, &cErr):
		c.JSON(cErr.Code, cErr)
	default:
		c.NoContent(http.StatusInternalServerError)
	}
}
