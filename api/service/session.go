package service

import (
	"context"
	"net/http"
	"time"

	"github.com/heiytor/invenda/api/pkg/cache"
	"github.com/heiytor/invenda/api/pkg/errors"
	"github.com/heiytor/invenda/api/pkg/hash"
	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/heiytor/invenda/api/pkg/requests"
)

type Session interface {
	CreateSession(ctx context.Context, req *requests.CreateSession) (insertedID string, err error)
	UpdateSession(ctx context.Context, req *requests.UpdateSession) (err error)
}

func (s *service) CreateSession(ctx context.Context, req *requests.CreateSession) (string, error) {
	usr, err := s.store.User.GetByEmail(ctx, req.Identifier)
	if err != nil {
		return "", errors.
			New().
			Attr("internal", err).
			Code(http.StatusNotFound).
			Layer(errors.LayerService).
			Msg("wrong identifer and/or password")
	}

	if !hash.Compare(req.Password, usr.Password) {
		return "", errors.
			New().
			Code(http.StatusNotFound).
			Layer(errors.LayerService).
			Msg("wrong identifer and/or password")
	}

	session := &models.Session{
		UserID:   usr.ID,
		State:    models.SessionStateActive,
		SourceIP: "TODO", // TODO
	}

	insertedID, err := s.store.Session.Create(ctx, session)

	// Associates the session with the user's preferred namespace and sets cache values.

	ns := new(models.Namespace)
	if usr.PreferredNamespace != "" {
		if ns, err = s.store.Namespace.Get(ctx, usr.PreferredNamespace); err != nil {
			return "", errors.
				New().
				Attr("internal", err).
				Code(http.StatusUnauthorized).
				Layer(errors.LayerService).
				Msg("user does not have any namespace")
		}
	} else {
		if ns, err = s.store.Namespace.GetFirst(ctx, usr.ID); err != nil {
			return "", errors.
				New().
				Attr("internal", err).
				Code(http.StatusUnauthorized).
				Layer(errors.LayerService).
				Msg("user does not have any namespace")
		}
	}

	member, _ := ns.FindMember(usr.ID)

	content := ns.ID + ";" + usr.ID + ";" + member.Permissions.String()
	if err := s.cache.Set(ctx, insertedID, content, cache.WithTTL(30*time.Minute)); err != nil {
		return "", err
	}

	return insertedID, nil
}

func (s *service) UpdateSession(ctx context.Context, req *requests.UpdateSession) error {
	ss, err := s.store.Session.Get(ctx, req.ID)
	if err != nil {
		return mapError(err, s.store.Session.Entity())
	}

	ns := new(models.Namespace)
	if req.Namespace != "" {
		if ns, err = s.store.Namespace.Get(ctx, req.Namespace); err != nil {
			return mapError(err, s.store.Namespace.Entity())
		}

		if err := s.store.User.Update(ctx, ss.UserID, &models.UserChanges{PreferredNamespace: ns.ID}); err != nil {
			return mapError(err, s.store.Namespace.Entity())
		}
	} else {
		if ns, err = s.store.Namespace.Get(ctx, ss.NamespaceID); err != nil {
			return mapError(err, s.store.Namespace.Entity())
		}
	}

	member, err := ns.FindMember(ss.UserID)
	if err != nil {
		return err
	}

	content := ns.ID + ";" + member.ID + ";" + member.Permissions.String()
	return s.cache.Set(ctx, req.ID, content, cache.WithTTL(30*time.Minute))
}
