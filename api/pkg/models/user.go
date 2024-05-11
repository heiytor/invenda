package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/heiytor/invenda/api/pkg/auth"
	"github.com/heiytor/invenda/api/pkg/clock"
)

type User struct {
	ID        string    `json:"id" bson:"_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password,omitempty" bson:"password"`

	// PreferredNamespace specifies the namespace the user should use when logging in.
	// The value must be updated whenever the user switches the session to a different namespace.
	PreferredNamespace string `json:"-" bson:"preferred_namespace"`
}

type UserChanges struct {
	Name               string    `bson:"name,omitempty"`
	Email              string    `bson:"email,omitempty"`
	Password           string    `bson:"password,omitempty"`
	UpdatedAt          time.Time `bson:"updated_at,omitempty"`
	PreferredNamespace string    `bson:"preferred_namespace,omitempty"`
}

type UserClaims struct {
	Email       string            `json:"email"`
	Namespace   string            `json:"namespace"`
	Permissions []auth.Permission `json:"permissions"`
	jwt.RegisteredClaims
}

func (u *UserClaims) SetRegisteredClaims() {
	u.RegisteredClaims.ID = uuid.New().String()
	u.RegisteredClaims.Issuer = "TODO"
	u.RegisteredClaims.Audience = jwt.ClaimStrings{"TODO"}

	now := clock.Now()
	u.RegisteredClaims.IssuedAt = jwt.NewNumericDate(now)
	u.RegisteredClaims.NotBefore = jwt.NewNumericDate(now)
	u.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(now.AddDate(0, 0, 7))
}
