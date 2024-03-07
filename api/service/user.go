package service

import (
	"context"
	"time"

	"github.com/heiytor/invenda/api/pkg/errors"
	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/heiytor/invenda/api/pkg/requests"
)

type User interface {
	UserCreate(ctx context.Context, req *requests.UserCreate) (string, error)
}

func (s *service) UserCreate(ctx context.Context, req *requests.UserCreate) (string, *errors.Error) {
	if conflicts, _ := s.store.UserConflicts(ctx, map[string]string{"email": req.Email}); len(conflicts) > 0 {
		return "", errors.New().Code(422).Attr("conflicts", "email").Msg("")
	}

	user := &models.User{
		Name:      req.Name,
		Email:     req.Email,
		Confirmed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	insertedID, err := s.store.UserCreate(ctx, user)
	if err != nil {
		panic(err)
	}

	return insertedID, nil
}
