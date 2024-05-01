package service

import (
	"context"
	"net/http"

	"github.com/heiytor/invenda/api/pkg/auth"
	"github.com/heiytor/invenda/api/pkg/errors"
	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/heiytor/invenda/api/pkg/requests"
	"github.com/heiytor/invenda/api/store"
)

type Namespace interface {
	ListNamespace(ctx context.Context, userID string, req *requests.ListNamespace) (nss []models.Namespace, count int64, err error)
	GetNamespace(ctx context.Context, memberID, namespaceID string) (ns *models.Namespace, err error)
	CreateNamespace(ctx context.Context, ownerID string, req *requests.CreateNamespace) (insertedID string, err error)
	UpdateNamespace(ctx context.Context, memberID, namespaceID string, req *requests.UpdateNamespace) (err error)
	DeleteNamespace(ctx context.Context, namespaceID string) (err error)
}

func (s *service) ListNamespace(ctx context.Context, userID string, req *requests.ListNamespace) ([]models.Namespace, int64, error) {
	nss, count, err := s.store.Namespace.GetMany(ctx, userID, &req.Query, store.ShortNamespace())
	return nss, count, mapError(err, s.store.Namespace.Entity())
}

func (s *service) GetNamespace(ctx context.Context, memberID, namespaceID string) (*models.Namespace, error) {
	ns, err := s.store.Namespace.Get(ctx, namespaceID)
	if err != nil {
		return nil, mapError(err, s.store.Namespace.Entity())
	}

	member, err := ns.FindMember(memberID)
	if err != nil {
		return nil, err
	}

	if !member.Owner {
		ns.WithoutPermissions()
	}

	return ns, nil
}

func (s *service) CreateNamespace(ctx context.Context, ownerID string, req *requests.CreateNamespace) (string, error) {
	if _, err := s.store.User.GetByID(ctx, ownerID); err != nil {
		return "", mapError(err, s.store.User.Entity())
	}

	for _, m := range req.Members {
		if _, err := s.store.User.GetByID(ctx, m.ID); err != nil {
			return "", mapError(err, s.store.User.Entity())
		}
	}

	ns := &models.Namespace{
		Name: req.Name,
		Members: append([]models.Member{
			{
				ID:          ownerID,
				Owner:       true,
				Permissions: auth.All(),
			},
		}, req.Members...),
	}

	insertedID, err := s.store.Namespace.Create(ctx, ns)
	return insertedID, mapError(err, s.store.Namespace.Entity())
}

func (s *service) UpdateNamespace(ctx context.Context, memberID, namespaceID string, req *requests.UpdateNamespace) error {
	// TODO: use transactions here

	ns, err := s.store.Namespace.Get(ctx, namespaceID)
	if err != nil {
		return mapError(err, s.store.Namespace.Entity())
	}

	if _, err := ns.FindMember(memberID); err != nil {
		return err
	}

	changes := &models.NamespaceChanges{
		Name: req.Name,
	}

	if err := s.store.Namespace.Update(ctx, namespaceID, changes); err != nil {
		return mapError(err, s.store.Namespace.Entity())
	}

	for _, usr := range req.Members {
		switch usr.Operation {
		case "upsert":
			if _, err := s.store.User.GetByID(ctx, usr.ID); err != nil {
				return mapError(err, s.store.User.Entity())
			}

			if m, _ := ns.FindMember(usr.ID); m != nil && m.Owner {
				return errors.
					New().
					Layer(errors.LayerService).
					Attr("operation", usr.Operation).
					Code(http.StatusForbidden).
					Msg("update namespace's owner is not allowed")
			}

			member := &models.Member{
				ID:          usr.ID,
				Owner:       false,
				Permissions: usr.Permissions,
			}

			// TODO: notify member
			if err := s.store.Namespace.UpsertMember(ctx, namespaceID, member); err != nil {
				return mapError(err, s.store.Namespace.Entity())
			}
		case "remove":
			member, err := ns.FindMember(usr.ID)
			if err != nil {
				return err
			}

			// TODO: the owner can remove itself
			if member.Owner {
				return errors.
					New().
					Layer(errors.LayerService).
					Attr("operation", usr.Operation).
					Code(http.StatusForbidden).
					Msg("remove namespace's owner is not allowed")
			}

			if err := s.store.Namespace.RemoveMember(ctx, namespaceID, usr.ID); err != nil {
				return mapError(err, s.store.Namespace.Entity())
			}
		default:
			return errors.
				New().
				Layer(errors.LayerService).
				Attr("operation", usr.Operation).
				Code(http.StatusForbidden).
				Msg("invalid operation")
		}
	}

	return nil
}

func (s *service) DeleteNamespace(ctx context.Context, namespaceID string) error {
	// TODO: notify all members
	err := s.store.Namespace.Delete(ctx, namespaceID)
	return mapError(err, s.store.Namespace.Entity())
}
