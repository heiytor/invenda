package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/heiytor/invenda/api/pkg/secretkeys"
)

type Claims interface {
	jwt.Claims
	SetRegisteredClaims()
}

// Encode encodes a claims to a signed JWT.
func Encode[T Claims](claims T) string {
	claims.SetRegisteredClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	str, _ := token.SignedString(secretkeys.PrivateKey)
	return str
}

// Decode decodes the raw JWT to claims.
func Decode[T Claims](raw string, claims T) error {
	_, err := jwt.ParseWithClaims(raw, claims, eval, jwt.WithValidMethods([]string{"EdDSA"}))
	return err
}

// eval evaluates if a token t is a valid token.
func eval(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
		return nil, fmt.Errorf("unexpected signature method: %v", t.Header["alg"])
	}

	return secretkeys.PublicKey, nil
}
