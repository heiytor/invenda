package errors

import (
	goerrors "errors"
	"fmt"
)

func Is(err error, target error) bool {
	return goerrors.Is(err, target)
}

type Error struct {
	Msg   string                 `json:"error"`
	Layer Layer                  `json:"layer"`
	Attrs map[string]interface{} `json:"details"`
	Code  int                    `json:"-"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("code=%d,attributes={%+v},msg=%s", err.Code, err.Attrs, err.Msg)
}

type errorBuilder struct {
	code  int
	layer Layer
	attrs map[string]interface{}
	msg   string
}

// As converts a generic err [error] to [Error]. It returns nil if err is not of type
// [Error].
func As(err error) *Error {
	e := new(Error)
	if !goerrors.As(err, &e) {
		return nil
	}

	return e
}

// New creates a new [Error] with code 500.
func New() *errorBuilder {
	return &errorBuilder{
		code:  500,
		attrs: map[string]interface{}{},
		msg:   "",
	}
}

// Code sets error's code.
func (err *errorBuilder) Code(code int) *errorBuilder {
	err.code = code
	return err
}

// Layer sets error's layer.
func (err *errorBuilder) Layer(layer Layer) *errorBuilder {
	err.layer = layer
	return err
}

// Code sets an error's attribute with key k and val v.
func (err *errorBuilder) Attr(k string, v interface{}) *errorBuilder {
	err.attrs[k] = v
	return err
}

// Msg sets the error's message. It ends the builder stage, returning an [Error].
func (err *errorBuilder) Msg(msg string) *Error {
	err.msg = msg
	return &Error{
		Code:  err.code,
		Layer: err.layer,
		Attrs: err.attrs,
		Msg:   err.msg,
	}
}
