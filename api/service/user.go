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
	GetUser(ctx context.Context, req *requests.GetUser) (usr *models.User, err error)
	CreateUser(ctx context.Context, req *requests.CreateUser) (insertedID string, err error)
	UpdateUser(ctx context.Context, id string, req *requests.UpdateUser) (usr *models.User, err error)
	DeleteUser(ctx context.Context, id string) (err error)
	AuthUser(ctx context.Context, req *requests.AuthUser) (claims *models.UserClaims, token string, err error)
}

func (s *service) GetUser(ctx context.Context, req *requests.GetUser) (*models.User, error) {
	usr, err := s.store.User.GetByID(ctx, req.ID, store.RemoveUserPassword)
	return usr, mapError(err, s.store.User.Entity())
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
	return insertedID, mapError(err, s.store.User.Entity())
}

func (s *service) UpdateUser(ctx context.Context, id string, req *requests.UpdateUser) (*models.User, error) {
	req.Email = strings.ToLower(req.Email)

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
		return nil, mapError(err, s.store.User.Entity())
	}

	usr, err := s.store.User.GetByID(ctx, id, store.RemoveUserPassword)
	return usr, mapError(err, s.store.User.Entity())
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	return mapError(s.store.User.Delete(ctx, id), s.store.User.Entity())
}

func (s *service) AuthUser(ctx context.Context, req *requests.AuthUser) (*models.UserClaims, string, error) {
	usr, _ := s.store.User.GetByEmail(ctx, req.Identifier)
	if usr == nil {
		return nil, "", errors.
			New().
			Code(http.StatusNotFound).
			Layer(errors.LayerService).
			Msg("wrong identifer and/or password")
	}

	if !hash.Compare(req.Password, usr.Password) {
		return nil, "", errors.
			New().
			Code(http.StatusNotFound).
			Layer(errors.LayerService).
			Msg("wrong identifer and/or password")
	}

	claims := &models.UserClaims{Email: usr.Email}
	claims.RegisteredClaims.Subject = usr.ID

	return claims, jwt.Encode(claims), nil
}
