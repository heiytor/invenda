package validator

import (
	"net/http"
	"unicode"

	"github.com/gookit/validate"
	"github.com/heiytor/invenda/api/pkg/errors"
)

type Validator struct{}

func New() *Validator {
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
	})

	validate.AddValidator("password", IsPassword)

	validate.AddGlobalMessages(map[string]string{
		"password": "{field} must be between 8 and 64 characters long, and contain at least one number, one uppercase letter, one lowercase letter, and one special character",
	})

	return &Validator{}
}

func (*Validator) Validate(s interface{}) error {
	v := validate.Struct(s)
	v.Validate()

	if v.Errors.Empty() {
		return nil
	}

	err := errors.New().Code(http.StatusBadRequest)
	for name, errs := range v.Errors.All() {
		details := make([]string, 0)
		for _, detail := range errs {
			details = append(details, detail)
		}
		err.Attr(name, details)
	}

	return err.Layer(errors.LayerPkg).Msg("bad request")
}

// IsPassword reports wheter a pwd is a valid password.
// A valid password must be between 8 and 64 characters long, and contain at least one number,
// one uppercase letter, one lowercase letter, and one special character",
func IsPassword(pwd string) bool {
	// golang's regexp does not supports backtracking, which means that is impossible
	// to make a regex to matches with these rules.
	var number, upper, lower, special bool
	var letters int

	for _, c := range pwd {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsLower(c):
			lower = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		}

		letters++
	}

	return (letters >= 8 && letters <= 64) && (number && upper && lower && special)
}
