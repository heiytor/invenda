package store

import (
	"context"

	"github.com/heiytor/invenda/api/pkg/clock"
	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/heiytor/invenda/api/pkg/query"
	"github.com/heiytor/invenda/api/store/internal"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetNamespaceOption is a function to be evaluated when retrieving a namespace document to modify its content.
type GetNamespaceOption func(u *models.Namespace) error

// ShortNamespace returns a function that modifies a namespace by setting some predefined attributes to their zero values.
func ShortNamespace() GetNamespaceOption {
	return func(u *models.Namespace) error {
		u.Members = []models.Member{}

		return nil
	}
}

type Namespace interface {
	Entity

	// Get retrieves a namespace with the specified ID. It returns the namespace or an error if any.
	Get(ctx context.Context, id string, opts ...GetNamespaceOption) (namespace *models.Namespace, err error)

	// GetFirst retrieves the first namespace for the provided member ID. It returs the namespace or an error if any.
	GetFirst(ctx context.Context, memberID string, opts ...GetNamespaceOption) (namespace *models.Namespace, err error)

	// GetMany retrieves a list of namespaces where a user is a member. The set of options will be applied to each
	// retrieved namespace. It returns the list of namespace, the total count of the existent document and an error if any.
	GetMany(ctx context.Context, userID string, query *query.Query, opts ...GetNamespaceOption) (namespaces []models.Namespace, count int64, err error)

	// Create creates a new namespace with the provided data. It returns the inserted ID or an error
	// if any.
	Create(ctx context.Context, ns *models.Namespace) (insertedID string, err error)

	// Update updates a namespace with the specified changes and ID. It returns [ErrNotFound] if no namespace is found.
	Update(ctx context.Context, id string, changes *models.NamespaceChanges) (err error)

	// Delete deletes a namespace with the specified ID. It returns [ErrNotFound] if no namespace is found.
	Delete(ctx context.Context, id string) (err error)

	// UpsertMember upserts a member within a specified namespace using the given member ID. This method should be used
	// whenever you want to add or modify a member. It returns [ErrNotFound] if no namespace is found.
	UpsertMember(ctx context.Context, id string, member *models.Member) (err error)

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

func (n *namespace) GetFirst(ctx context.Context, memberID string, opts ...GetNamespaceOption) (*models.Namespace, error) {
	namespace := new(models.Namespace)
	if err := n.c.FindOne(ctx, bson.M{"members": bson.M{"$elemMatch": bson.M{"_id": memberID}}}).Decode(namespace); err != nil {
		return nil, mapError(err)
	}

	for _, opt := range opts {
		if err := opt(namespace); err != nil {
			return nil, err
		}
	}

	return namespace, nil
}

func (n *namespace) GetMany(ctx context.Context, userID string, query *query.Query, opts ...GetNamespaceOption) ([]models.Namespace, int64, error) {
	match := bson.M{"members": bson.M{"$elemMatch": bson.M{"_id": userID}}}

	count, err := n.c.CountDocuments(ctx, match)
	if err != nil {
		log.Error().Err(err).Msg("unable to count the total documents")
	}

	pipeline := make([]bson.M, 0)
	pipeline = append(pipeline, bson.M{"$match": match})
	pipeline = append(pipeline, internal.FromPaginator(&query.Paginator)...)
	pipeline = append(pipeline, internal.FromSorter(&query.Sorter)...)

	cursor, err := n.c.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, mapError(err)
	}
	defer cursor.Close(ctx)

	namespaces := make([]models.Namespace, 0)
	for cursor.Next(ctx) {
		ns := new(models.Namespace)
		if err := cursor.Decode(ns); err != nil {
			return nil, 0, mapError(err)
		}

		for _, opt := range opts {
			if err := opt(ns); err != nil {
				return nil, 0, err
			}
		}

		namespaces = append(namespaces, *ns)
	}

	return namespaces, count, nil
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

func (n *namespace) UpsertMember(ctx context.Context, id string, member *models.Member) error {
	filter := bson.M{"_id": id, "members": bson.M{"$elemMatch": bson.M{"_id": member.ID}}}
	upsert := bson.M{}

	// Add the member to the set if it does not already exist. Otherwise, upsert.
	if ok := n.c.FindOne(ctx, filter); ok.Err() != nil {
		member.AddedAt = clock.Now()
		filter = bson.M{"_id": id}
		upsert = bson.M{"$addToSet": bson.M{"members": member}}
	} else {
		upsert = bson.M{"$set": bson.M{"members.$": member}}
	}

	res, err := n.c.UpdateOne(ctx, filter, upsert)
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
