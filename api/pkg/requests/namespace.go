package requests

import (
	"github.com/heiytor/invenda/api/pkg/auth"
	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/heiytor/invenda/api/pkg/query"
)

type member struct {
	Operation   string            `json:"operation" validate:"in:upsert,remove|required"`
	ID          string            `json:"id" validate:"ulid|required"`
	Permissions []auth.Permission `json:"permissions" validate:"permissions|required"`
}

type ListNamespace struct {
	query.Query
}

type CreateNamespace struct {
	Name    string          `json:"name" validate:"required"`
	Members []models.Member `json:"members" validate:"permissions"`
}

type UpdateNamespace struct {
	Name    string   `json:"name" validate:""`
	Members []member `json:"members"`
}
