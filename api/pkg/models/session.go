package models

import (
	"time"

	"github.com/heiytor/invenda/api/pkg/auth"
)

type SessionState string

const (
	SessionStateActive   SessionState = "active"
	SessionStateInactive SessionState = "inactive"
)

type Session struct {
	ID          string           `json:"id" bson:"_id"`
	UserID      string           `json:"user_id" bson:"user_id"`
	SourceIP    string           `json:"source_ip" bson:"source_ip"`
	StartedAt   time.Time        `json:"started_at" bson:"started_at"`
	EndedAt     time.Time        `json:"ended_at" bson:"ended_at"`
	State       SessionState     `json:"state" bson:"state"`
	NamespaceID string           `json:"-" bson:"-"`
	Permissions auth.Permissions `json:"-" bson:"-"`
}

type SessionChanges struct {
	EndedAt time.Time    `bson:"ended_at,omitempty"`
	State   SessionState `bson:"state,omitempty"`
}
