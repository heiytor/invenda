package errors

import "fmt"

type Error struct {
	Code  int16
	Attrs map[string]string
	Msg   string
}

func (err *Error) Error() string {
	return fmt.Sprintf("%s code=%d attributes=%v", err.Msg, err.Code, err.Attrs)
}

type ErrorBuilder struct {
	code  int16
	attrs map[string]string
	msg   string
}

func New() *ErrorBuilder {
	return nil
}

func (err *ErrorBuilder) Code(code int16) *ErrorBuilder {
	err.code = code
	return err
}

func (err *ErrorBuilder) Attr(key, val string) *ErrorBuilder {
	err.attrs[key] = val
	return err
}

func (err *ErrorBuilder) Msg(msg string) *Error {
	err.msg = msg
	return &Error{
		Code:  err.code,
		Attrs: err.attrs,
		Msg:   err.msg,
	}
}
