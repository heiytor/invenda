package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/heiytor/invenda/api/store"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUserCreate(t *testing.T) {
	type Actual struct {
		err error
	}

	cases := []struct {
		description string
		user        *models.User
		expected    Actual
	}{
		{
			description: "succeeds to create a user",
			user: &models.User{
				Name:     "John Doe",
				Email:    "john.doe@test.com",
				Password: "secret",
			},
			expected: Actual{err: nil},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			defer srv.reset()
			ctx := context.Background()

			id, err := s.User.Create(ctx, tc.user)
			require.Equal(t, tc.expected, Actual{err})
			require.NotEmpty(t, id)

			user := new(models.User)
			require.NoError(t, db.Collection("user").FindOne(ctx, bson.M{"_id": id}).Decode(user))
			require.NotZero(t, user)
			require.Equal(t, tc.user.Name, user.Name)
			require.Equal(t, tc.user.Email, user.Email)
			require.Equal(t, tc.user.Password, user.Password)
		})
	}
}

func TestUserGetByID(t *testing.T) {
	type Actual struct {
		user *models.User
		err  error
	}

	cases := []struct {
		description string
		id          string
		opts        []store.GetUserOption
		fixtures    []fixture
		expected    Actual
	}{
		{
			description: "fails when user is not found",
			id:          "00000000000000000000000000",
			opts:        []store.GetUserOption{},
			fixtures:    []fixture{},
			expected: Actual{
				user: nil,
				err:  store.ErrNotFound,
			},
		},
		{
			description: "succeeds to find a user",
			id:          "01HNGJ2BTGQAHAZ1XNYZQPG719",
			opts:        []store.GetUserOption{},
			fixtures:    []fixture{fixtureUser},
			expected: Actual{
				user: &models.User{
					ID:        "01HNGJ2BTGQAHAZ1XNYZQPG719",
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					Name:      "John Doe",
					Email:     "john.doe@test.com",
					Password:  "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0",
				},
				err: nil,
			},
		},
		{
			description: "succeeds to find a user and remove the password",
			id:          "01HNGJ2BTGQAHAZ1XNYZQPG719",
			opts:        []store.GetUserOption{store.RemoveUserPassword},
			fixtures:    []fixture{fixtureUser},
			expected: Actual{
				user: &models.User{
					ID:        "01HNGJ2BTGQAHAZ1XNYZQPG719",
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					Name:      "John Doe",
					Email:     "john.doe@test.com",
				},
				err: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			srv.apply(tc.fixtures...)
			defer srv.reset()

			ctx := context.Background()

			user, err := s.User.GetByID(ctx, tc.id, tc.opts...)
			require.Equal(t, tc.expected, Actual{user, err})
		})
	}
}

func TestUserGetByEmail(t *testing.T) {
	type Actual struct {
		user *models.User
		err  error
	}

	cases := []struct {
		description string
		email       string
		opts        []store.GetUserOption
		fixtures    []fixture
		expected    Actual
	}{
		{
			description: "fails when user is not found",
			email:       "user@test.com",
			opts:        []store.GetUserOption{},
			fixtures:    []fixture{},
			expected: Actual{
				user: nil,
				err:  store.ErrNotFound,
			},
		},
		{
			description: "succeeds to find a user",
			email:       "john.doe@test.com",
			opts:        []store.GetUserOption{},
			fixtures:    []fixture{fixtureUser},
			expected: Actual{
				user: &models.User{
					ID:        "01HNGJ2BTGQAHAZ1XNYZQPG719",
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					Name:      "John Doe",
					Email:     "john.doe@test.com",
					Password:  "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0",
				},
				err: nil,
			},
		},
		{
			description: "succeeds to find a user and remove the password",
			email:       "john.doe@test.com",
			opts:        []store.GetUserOption{store.RemoveUserPassword},
			fixtures:    []fixture{fixtureUser},
			expected: Actual{
				user: &models.User{
					ID:        "01HNGJ2BTGQAHAZ1XNYZQPG719",
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					Name:      "John Doe",
					Email:     "john.doe@test.com",
				},
				err: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			srv.apply(tc.fixtures...)
			defer srv.reset()

			ctx := context.Background()

			user, err := s.User.GetByEmail(ctx, tc.email, tc.opts...)
			require.Equal(t, tc.expected, Actual{user, err})
		})
	}
}

func TestUserConflicts(t *testing.T) {
	type Actual struct {
		conflicts []string
		err       error
	}

	cases := []struct {
		description string
		target      *models.User
		fixtures    []fixture
		expected    Actual
	}{
		{
			description: "succeeds when none conflicts are found",
			target: &models.User{
				Email: "john.doe@test.com",
			},
			fixtures: []fixture{},
			expected: Actual{
				conflicts: []string{},
				err:       nil,
			},
		},
		{
			description: "succeeds when a conflict is found",
			target: &models.User{
				Email: "john.doe@test.com",
			},
			fixtures: []fixture{fixtureUser},
			expected: Actual{
				conflicts: []string{"email"},
				err:       nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			srv.apply(tc.fixtures...)
			defer srv.reset()

			ctx := context.Background()

			conflicts, err := s.User.Conflicts(ctx, tc.target)
			require.Equal(t, tc.expected, Actual{conflicts, err})
		})
	}
}

func TestUserUpdate(t *testing.T) {
	type Actual struct {
		err error
	}

	cases := []struct {
		description string
		id          string
		changes     *models.UserChanges
		fixtures    []fixture
		expected    Actual
	}{
		{
			description: "fails when user is not found",
			id:          "00000000000000000000000000",
			changes:     &models.UserChanges{},
			fixtures:    []fixture{},
			expected:    Actual{err: store.ErrNotFound},
		},
		{
			description: "succeeds to find a user",
			id:          "01HNGJ2BTGQAHAZ1XNYZQPG719",
			changes:     &models.UserChanges{Email: "new.email@test.com"},
			fixtures:    []fixture{fixtureUser},
			expected:    Actual{err: nil},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			srv.apply(tc.fixtures...)
			defer srv.reset()

			ctx := context.Background()

			if err := s.User.Update(ctx, tc.id, tc.changes); err != nil {
				require.Equal(t, tc.expected, Actual{err})
				return
			}

			user := new(models.User)
			require.NoError(t, db.Collection("user").FindOne(ctx, bson.M{"_id": tc.id}).Decode(user))
			require.NotZero(t, user)

			require.Equal(t, tc.changes.Email, user.Email)                                    // Checks if the email was updated
			require.WithinDuration(t, tc.changes.UpdatedAt, user.UpdatedAt, time.Millisecond) // Checks if the updated_at was updated
			require.Equal(t, "John Doe", user.Name)                                           // Checks if the password follows unalterated
		})
	}
}
