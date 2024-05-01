package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/heiytor/invenda/api/pkg/clock"
)

type User struct {
	ID        string    `json:"id" bson:"_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password,omitempty" bson:"password"`
}

type UserChanges struct {
	Name      string    `bson:"name,omitempty"`
	Email     string    `bson:"email,omitempty"`
	Password  string    `bson:"password,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
}

type UserClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (u *UserClaims) SetRegisteredClaims() {
	u.RegisteredClaims.ID = "TODO UUID"
	u.RegisteredClaims.Issuer = "TODO"
	u.RegisteredClaims.Audience = jwt.ClaimStrings{"TODO"}
	u.RegisteredClaims.IssuedAt = jwt.NewNumericDate(clock.Now())
	u.RegisteredClaims.NotBefore = jwt.NewNumericDate(clock.Now())
	u.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(clock.Now().AddDate(0, 0, 7)) // TODO: maybe we can expires after 1 week?)
}
