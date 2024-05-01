package models

import "time"

type Namespace struct {
	ID        string    `json:"id" bson:"_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Name      string    `json:"name" bson:"name"`
	Members   []Member  `json:"members" bson:"members"`
}

type Member struct {
	ID          string    `json:"id" bson:"_id"`
	AddedAt     time.Time `json:"created_at" bson:"added_at"`
	Permissions []string  `json:"permissions" bson:"permissions"` // TODO: permissions
}

type NamespaceChanges struct {
	UpdatedAt time.Time `bson:"updated_at"`
	Name      string    `bson:"name,omitempty"`
}
