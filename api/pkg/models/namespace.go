package models

import (
	"net/http"
	"slices"
	"time"

	"github.com/heiytor/invenda/api/pkg/auth"
	"github.com/heiytor/invenda/api/pkg/errors"
)

type Namespace struct {
	ID        string    `json:"id" bson:"_id"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
	Name      string    `json:"name" bson:"name"`
	Members   []Member  `json:"members,omitempty" bson:"members"`
}

// FindMember reports whether a member exists or not in the namespace.
func (ns *Namespace) FindMember(id string) (*Member, error) {
	member := new(Member)
	ok := false

	containsFunc := func(m Member) bool {
		if ok = m.ID == id; ok {
			member = &m
		}

		return ok
	}

	if !slices.ContainsFunc(ns.Members, containsFunc) {
		return nil, errors.New().Attr("member_id", id).Code(http.StatusForbidden).Msg(errors.MsgMemberNotFound)
	}

	return member, nil
}

// WithoutPermissions clears the permissions of all members in the namespace.
// This method can used to restrict view when a member does not have correct permissions.
func (ns *Namespace) WithoutPermissions() {
	for i := range ns.Members {
		ns.Members[i].Permissions = []auth.Permission{}
	}
}

type Member struct {
	ID          string           `json:"id" bson:"_id"`
	AddedAt     time.Time        `json:"added_at" bson:"added_at"`
	Owner       bool             `json:"owner" bson:"owner"`
	Permissions auth.Permissions `json:"permissions,omitempty" bson:"permissions"`
}

type NamespaceChanges struct {
	UpdatedAt time.Time `bson:"updated_at"`
	Name      string    `bson:"name,omitempty"`
}
