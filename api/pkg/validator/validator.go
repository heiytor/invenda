package validator

import (
	"net/http"
	"slices"
	"strings"
	"unicode"

	"github.com/gookit/validate"
	"github.com/heiytor/invenda/api/pkg/auth"
	"github.com/heiytor/invenda/api/pkg/errors"
	"github.com/oklog/ulid/v2"
)

type Validator struct{}

func New() *Validator {
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
	})

	validate.AddValidator("password", IsPassword)
	validate.AddValidator("ulid", IsULID)
	validate.AddValidator("permissions", IsPermissions)

	validate.AddGlobalMessages(map[string]string{
		"password":    "{field} must be between 8 and 64 characters long, and contain at least one number, one uppercase letter, one lowercase letter, and one special character.",
		"ulid":        "{field} must be a valid ULID.",
		"permissions": "{field} must be a valid permission",
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

// IsPassword reports wheter a input is a valid password.
// A valid password must be between 8 and 64 characters long, and contain at least one number,
// one uppercase letter, one lowercase letter, and one special character",
func IsPassword(input string) bool {
	// golang's regexp does not supports backtracking, which means that is impossible
	// to make a regex to matches with these rules.
	var number, upper, lower, special bool
	var letters int

	for _, c := range input {
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

func IsULID(input string) bool {
	parts := strings.Split(input, "_")

	if len(parts) != 2 {
		return false
	}

	if _, err := ulid.Parse(parts[1]); err != nil {
		return false
	}

	return true
}

func IsPermissions(inputs []auth.Permission) bool {
	all := auth.All()
	for _, input := range inputs {
		if !slices.Contains(all, input) {
			return false
		}
	}

	return true
}
