package store

import (
	"context"

	"github.com/heiytor/invenda/api/pkg/clock"
	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUserOption is a function to be evaluated when retrieving a user document to modify its content.
type GetUserOption func(u *models.User) error

// RemoveUserPassword sets the user's password to an empty string
func RemoveUserPassword(u *models.User) error {
	u.Password = ""
	return nil
}

type User interface {
	Entity

	// Create creates a new user with the provided data. It returns the inserted ID or an error
	// if any.
	Create(ctx context.Context, user *models.User) (insertedID string, err error)

	// GetByID retrieves a user with the specified ID. It returns the found user or an error if any.
	GetByID(ctx context.Context, id string, opts ...GetUserOption) (user *models.User, err error)

	// GetByEmail retrieves a user with the specified email. It returns the found user or an error if any.
	GetByEmail(ctx context.Context, email string, opts ...GetUserOption) (user *models.User, err error)

	// Conflicts reports whether the fields of the provided target already exist in the database.
	// It returns a list of conflicted fields or an error if any.
	Conflicts(ctx context.Context, target *models.User) (conflicts []string, err error)

	// Update modifies a user with the specified ID based on the provided changes. It returns [ErrNotFound]
	// if no user is found.
	Update(ctx context.Context, id string, changes *models.UserChanges) (err error)

	// Delete deletes a user with the specified ID. It returns [ErrNotFound] if no user is found.
	Delete(ctx context.Context, id string) (err error)
}

type user struct {
	c *mongo.Collection // c is the "user" collection
}

var _ User = (*user)(nil)

func (*user) Entity() string {
	return "user"
}

func (u *user) Create(ctx context.Context, usr *models.User) (string, error) {
	usr.ID = "usr_" + ulid.Make().String()
	usr.CreatedAt = clock.Now()
	usr.UpdatedAt = clock.Now()

	if _, err := u.c.InsertOne(ctx, usr); err != nil {
		return "", mapError(err)
	}

	return usr.ID, nil
}

func (u *user) GetByID(ctx context.Context, id string, opts ...GetUserOption) (*models.User, error) {
	usr := new(models.User)
	if err := u.c.FindOne(ctx, bson.M{"_id": id}).Decode(usr); err != nil {
		return nil, mapError(err)
	}

	for _, opt := range opts {
		if err := opt(usr); err != nil {
			return nil, err
		}
	}

	return usr, nil
}

func (u *user) GetByEmail(ctx context.Context, email string, opts ...GetUserOption) (*models.User, error) {
	usr := new(models.User)
	if err := u.c.FindOne(ctx, bson.M{"email": email}).Decode(usr); err != nil {
		return nil, mapError(err)
	}

	for _, opt := range opts {
		if err := opt(usr); err != nil {
			return nil, err
		}
	}

	return usr, nil
}

func (u *user) Conflicts(ctx context.Context, target *models.User) ([]string, error) {
	cursor, err := u.c.Aggregate(ctx, or(target))
	if err != nil {
		return nil, mapError(err)
	}
	defer cursor.Close(ctx)

	conflicts := make([]string, 0)
	for cursor.Next(ctx) {
		usr := new(models.User)

		if err = cursor.Decode(usr); err != nil {
			return nil, mapError(err)
		}

		conflicts = append(conflicts, partialEqual(target, usr)...)
	}

	return conflicts, nil
}

func (u *user) Update(ctx context.Context, id string, changes *models.UserChanges) error {
	if changes == nil {
		return nil
	}

	changes.UpdatedAt = clock.Now()

	res, err := u.c.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": changes})
	if err != nil {
		return mapError(err)
	}

	if res.MatchedCount < 1 {
		return ErrNotFound
	}

	return nil
}

func (u *user) Delete(ctx context.Context, id string) error {
	res, err := u.c.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return mapError(err)
	}

	if res.DeletedCount < 1 {
		return ErrNotFound
	}

	return nil
}
