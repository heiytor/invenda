package store

import (
	"context"

	"github.com/heiytor/invenda/api/pkg/clock"
	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetNamespaceOption TODO
type GetNamespaceOption func(u *models.Namespace) error

type Namespace interface {
	Entity

	// Get retrieves a namespace with the specified ID. It returns the namespace or an error if any.
	Get(ctx context.Context, id string, opts ...GetNamespaceOption) (namespace *models.Namespace, err error)

	// Create creates a new namespace with the provided data. It returns the inserted ID or an error
	// if any.
	Create(ctx context.Context, ns *models.Namespace) (insertedID string, err error)

	// Update updates a namespace with the specified changes and ID. It returns [ErrNotFound] if no namespace is found.
	Update(ctx context.Context, id string, changes *models.NamespaceChanges) (err error)

	// Delete deletes a namespace with the specified ID. It returns [ErrNotFound] if no namespace is found.
	Delete(ctx context.Context, id string) (err error)

	UpsertMember(ctx context.Context, id string, member *models.Member, exists bool) (err error)

	// RemoveMember removes a member with the specified memberID from the namespace with the specified id.
	// It returns [ErrNotFound] if no namespace or member is found.
	RemoveMember(ctx context.Context, id, memberID string) (err error)
}

type namespace struct {
	c *mongo.Collection // c is the "user" collection
}

var _ Namespace = (*namespace)(nil)

func (*namespace) Entity() string {
	return "namespace"
}

func (n *namespace) Get(ctx context.Context, id string, opts ...GetNamespaceOption) (*models.Namespace, error) {
	namespace := new(models.Namespace)
	if err := n.c.FindOne(ctx, bson.M{"_id": id}).Decode(namespace); err != nil {
		return nil, mapError(err)
	}

	for _, opt := range opts {
		if err := opt(namespace); err != nil {
			return nil, err
		}
	}

	return namespace, nil
}

func (n *namespace) Create(ctx context.Context, ns *models.Namespace) (string, error) {
	ns.ID = "ns_" + ulid.Make().String()

	now := clock.Now()
	ns.CreatedAt = now
	ns.UpdatedAt = now

	for _, m := range ns.Members {
		m.AddedAt = now
	}

	if _, err := n.c.InsertOne(ctx, ns); err != nil {
		return "", mapError(err)
	}

	return ns.ID, nil
}

func (n *namespace) Update(ctx context.Context, id string, changes *models.NamespaceChanges) error {
	if changes == nil {
		return nil
	}

	changes.UpdatedAt = clock.Now()

	res, err := n.c.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": changes})
	if err != nil {
		return mapError(err)
	}

	if res.MatchedCount < 1 {
		return ErrNotFound
	}

	return nil
}

func (n *namespace) Delete(ctx context.Context, id string) error {
	res, err := n.c.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return mapError(err)
	}

	if res.DeletedCount < 1 {
		return ErrNotFound
	}

	return nil
}

func (n *namespace) UpsertMember(ctx context.Context, id string, member *models.Member, exists bool) error {
	update := make(bson.M)
	if exists {
		update = bson.M{"$set": bson.M{"members.$": member}}
	} else {
		member.AddedAt = clock.Now()
		update = bson.M{"$addToSet": bson.M{"members": member}}
	}

	res, err := n.c.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return mapError(err)
	}

	if res.MatchedCount < 1 {
		return ErrNotFound
	}

	return nil
}

func (n *namespace) RemoveMember(ctx context.Context, id, memberID string) error {
	res, err := n.c.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$pull": bson.M{"members": bson.M{"_id": memberID}}})
	if err != nil {
		return mapError(err)
	}

	if res.MatchedCount < 1 {
		return ErrNotFound
	}

	return nil
}
