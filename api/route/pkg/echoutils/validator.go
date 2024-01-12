package echoutils

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	val *validator.Validate
}

func (v *Validator) Validate(s interface{}) error {
	if v.val == nil {
		v.val = validator.New()
	}

	if err := v.val.Struct(s); err != nil {
		// TODO: mapear os erros para o status code correto
		return errors.New("TODO")
	}

	return nil
}
