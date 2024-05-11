package store

import (
	"context"
	"time"

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
type GetSessionOption func(s *models.Session) error

type Session interface {
	Entity

	Get(ctx context.Context, id string, opts ...GetSessionOption) (session *models.Session, err error)
	List(ctx context.Context, userID string, query *query.Query, opts ...GetSessionOption) (sessions []models.Session, count int64, err error)
	Create(ctx context.Context, session *models.Session) (insertedID string, err error)
	Update(ctx context.Context, id string, changes *models.SessionChanges) (err error)
	Delete(ctx context.Context, id string) (err error)
}

type session struct {
	c *mongo.Collection // c is the "session" collection
}

var _ Session = (*session)(nil) // Ensures that session implements Session

func (*session) Entity() string {
	return "session"
}

func (s *session) Get(ctx context.Context, id string, opts ...GetSessionOption) (*models.Session, error) {
	ss := new(models.Session)
	if err := s.c.FindOne(ctx, bson.M{"_id": id}).Decode(ss); err != nil {
		return nil, mapError(err)
	}

	for _, opt := range opts {
		if err := opt(ss); err != nil {
			return nil, err
		}
	}

	return ss, nil
}

func (s *session) List(ctx context.Context, userID string, query *query.Query, opts ...GetSessionOption) ([]models.Session, int64, error) {
	match := bson.M{"userID": userID}

	count, err := s.c.CountDocuments(ctx, match)
	if err != nil {
		log.Error().Err(err).Msg("unable to count the total documents")
	}

	pipeline := make([]bson.M, 0)
	pipeline = append(pipeline, bson.M{"$match": match})
	pipeline = append(pipeline, internal.FromPaginator(&query.Paginator)...)
	pipeline = append(pipeline, internal.FromSorter(&query.Sorter)...)

	cursor, err := s.c.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, mapError(err)
	}
	defer cursor.Close(ctx)

	sessions := make([]models.Session, 0)
	for cursor.Next(ctx) {
		ss := new(models.Session)
		if err := cursor.Decode(ss); err != nil {
			return nil, 0, mapError(err)
		}

		for _, opt := range opts {
			if err := opt(ss); err != nil {
				return nil, 0, err
			}
		}

		sessions = append(sessions, *ss)
	}

	return sessions, count, nil
}

func (s *session) Create(ctx context.Context, ss *models.Session) (string, error) {
	ss.ID = "ss_" + ulid.Make().String()
	ss.StartedAt = clock.Now()
	ss.EndedAt = time.Time{}

	if _, err := s.c.InsertOne(ctx, ss); err != nil {
		return "", mapError(err)
	}

	return ss.ID, nil
}

func (s *session) Update(ctx context.Context, id string, changes *models.SessionChanges) error {
	if changes == nil {
		return nil
	}

	res, err := s.c.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": changes})
	if err != nil {
		return mapError(err)
	}

	if res.MatchedCount < 1 {
		return ErrNotFound
	}

	return nil
}

func (s *session) Delete(ctx context.Context, id string) error {
	res, err := s.c.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return mapError(err)
	}

	if res.DeletedCount < 1 {
		return ErrNotFound
	}

	return nil
}
