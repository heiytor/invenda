package service

import (
	"context"
	"net/http"
	"strings"

	"github.com/heiytor/invenda/api/pkg/errors"
	"github.com/heiytor/invenda/api/pkg/hash"
	"github.com/heiytor/invenda/api/pkg/jwt"
	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/heiytor/invenda/api/pkg/requests"
	"github.com/heiytor/invenda/api/store"
)

type User interface {
	CreateUser(ctx context.Context, req *requests.CreateUser) (insertedID string, err error)
	UpdateUser(ctx context.Context, id string, req *requests.UpdateUser) (usr *models.User, err error)
	AuthUser(ctx context.Context, req *requests.AuthUser) (token *models.UserClaims, err error)
}

func (s *service) CreateUser(ctx context.Context, req *requests.CreateUser) (string, error) {
	if conflicts, _ := s.store.User.Conflicts(ctx, &models.User{Email: req.Email}); len(conflicts) > 0 {
		return "", errors.
			New().
			Code(http.StatusConflict).
			Attr("entity", "user").
			Attr("conflicts", conflicts).
			Layer(errors.LayerService).
			Msg(errors.MsgConflict)
	}

	usr := &models.User{
		Name:     req.Name,
		Email:    strings.ToLower(req.Email),
		Password: hash.New(req.Password, hash.DefaultOptions),
	}

	insertedID, err := s.store.User.Create(ctx, usr)
	return insertedID, err // TODO: map this error
}

func (s *service) UpdateUser(ctx context.Context, id string, req *requests.UpdateUser) (*models.User, error) {
	if conflicts, _ := s.store.User.Conflicts(ctx, &models.User{Email: req.Email}); len(conflicts) > 0 {
		return nil, errors.
			New().
			Code(http.StatusConflict).
			Attr("entity", "user").
			Attr("conflicts", conflicts).
			Layer(errors.LayerService).
			Msg(errors.MsgConflict)
	}

	changes := &models.UserChanges{
		Name:  req.Name,
		Email: req.Email,
	}

	if req.Password != "" {
		changes.Password = hash.New(req.Password, hash.DefaultOptions)
	}

	if err := s.store.User.Update(ctx, id, changes); err != nil {
		return nil, err // TODO: map this error
	}

	usr, err := s.store.User.GetByID(ctx, id, store.RemoveUserPassword)
	return usr, err // TODO: map this error
}

func (s *service) AuthUser(ctx context.Context, req *requests.AuthUser) (*models.UserClaims, error) {
	usr, err := s.store.User.GetByEmail(ctx, req.Identifier)
	if err != nil {
		return nil, err // TODO: map err
	}

	if !hash.Compare(req.Password, usr.Password) {
		return nil, errors.New().Msg("foo") //TODO: map err
	}

	claims := &models.UserClaims{Email: usr.Email}
	claims.RegisteredClaims.Subject = usr.ID
	claims.Token = jwt.Encode(claims)

	return claims, nil
}
