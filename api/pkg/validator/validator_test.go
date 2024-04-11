package validator_test

import (
	"testing"

	"github.com/heiytor/invenda/api/pkg/validator"
	"github.com/stretchr/testify/require"
)

func TestIsPassword(t *testing.T) {
	cases := []struct {
		description string
		password    string
		expected    bool
	}{
		{
			description: "password with less than 8 characters",
			password:    "Pwd1$",
			expected:    false,
		},
		{
			description: "password with more than 64 characters",
			password:    "ThisIsAVeryLongPasswordThatExceeds64CharactersInLengthAndIsNotValid1$",
			expected:    false,
		},
		{
			description: "password without a number",
			password:    "password$",
			expected:    false,
		},
		{
			description: "password without a special character",
			password:    "Password1",
			expected:    false,
		},
		{
			description: "password without an uppercase letter",
			password:    "password1$",
			expected:    false,
		},
		{
			description: "password without a lowercase letter",
			password:    "PASSWORD1$",
			expected:    false,
		},
		{
			description: "password with only numbers",
			password:    "123456789",
			expected:    false,
		},
		{
			description: "password with only special characters",
			password:    "!@#$%^&*()",
			expected:    false,
		},
		{
			description: "password with only uppercase letters",
			password:    "PASSWORD",
			expected:    false,
		},
		{
			description: "password with only lowercase letters",
			password:    "password",
			expected:    false,
		},
		{
			description: "Valid password",
			password:    "Password1$",
			expected:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			require.Equal(t, tc.expected, validator.IsPassword(tc.password))
		})
	}
}
