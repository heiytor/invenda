package store

import (
	"context"
	"testing"
	"time"

	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/heiytor/invenda/api/store/dbtest"
	"github.com/stretchr/testify/assert"
)

func TestUserCreate(t *testing.T) {
	cases := []struct {
		description string
		user        *models.User
		expected    error
	}{
		{
			description: "succeeds to create a user",
			user: &models.User{
				Name:      "John Doe",
				Email:     "john_doe@test.com",
				Confirmed: true,
				CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			expected: nil,
		},
	}

	db, err := dbtest.New()
	if err != nil {
		t.Fatalf("Failed to create the test container: %s", err.Error())
	}
	defer db.Teardown()

	store := New(context.TODO(), db.URI)

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			_, err := store.UserCreate(context.TODO(), tc.user)
			assert.NoError(t, err)
		})
	}
}

func TestUserConflicts(t *testing.T) {
	type Expected struct {
		conflicts []string
		err       error
	}

	cases := []struct {
		description string
		fixtures    []string
		fields      map[string]string
		expected    Expected
	}{
		{
			description: "succeeds with an empty slice if no conflicts are found",
			fixtures:    []string{dbtest.FixtureUser},
			fields:      map[string]string{"email": "nonexistent"},
			expected: Expected{
				conflicts: []string{},
				err:       nil,
			},
		},
		{
			description: "succeeds if a conflict is found",
			fixtures:    []string{dbtest.FixtureUser},
			fields:      map[string]string{"email": "john_doe@test.com"},
			expected: Expected{
				conflicts: []string{"email"},
				err:       nil,
			},
		},
	}

	db, err := dbtest.New()
	if err != nil {
		t.Fatalf("Failed to create the test container: %s", err.Error())
	}
	defer db.Teardown()

	store := New(context.TODO(), db.URI)

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, db.ApplyFixtures(tc.fixtures...))
			defer db.TeardownFixtures()

			conflicts, err := store.UserConflicts(context.TODO(), tc.fields)
			assert.Equal(t, tc.expected, Expected{conflicts, err})
		})
	}
}

func TestUserGet(t *testing.T) {
	type Expected struct {
		user *models.User
		err  error
	}

	cases := []struct {
		description string
		fixtures    []string
		id          string
		expected    Expected
	}{
		{
			description: "fails when the user is not found",
			fixtures:    []string{dbtest.FixtureUser},
			id:          "00000000000000000000000000",
			expected: Expected{
				user: nil,
				err:  ErrNotFound,
			},
		},
		{
			description: "succeeds to retrieve a user",
			fixtures:    []string{dbtest.FixtureUser},
			id:          "01HNGJ2BTGQAHAZ1XNYZQPG719",
			expected: Expected{
				user: &models.User{
					ID:        "01HNGJ2BTGQAHAZ1XNYZQPG719",
					Name:      "John Doe",
					Email:     "john_doe@test.com",
					Confirmed: true,
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				},
				err: nil,
			},
		},
	}

	db, err := dbtest.New()
	if err != nil {
		t.Fatalf("Failed to create the test container: %s", err.Error())
	}
	defer db.Teardown()

	store := New(context.TODO(), db.URI)

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, db.ApplyFixtures(tc.fixtures...))
			defer db.TeardownFixtures()

			user, err := store.UserGet(context.TODO(), tc.id)
			assert.Equal(t, tc.expected, Expected{user, err})
		})
	}
}

func TestUserSetConfirmed(t *testing.T) {
	cases := []struct {
		description string
		fixtures    []string
		id          string
		confirmed   bool
		expected    error
	}{
		{
			description: "fails when the user is not found",
			fixtures:    []string{dbtest.FixtureUser},
			id:          "00000000000000000000000000",
			confirmed:   true,
			expected:    ErrNotFound,
		},
		{
			description: "succeeds to confirm a user",
			fixtures:    []string{dbtest.FixtureUser},
			id:          "01HNGJ2BTGQAHAZ1XNYZQPG719",
			confirmed:   true,
			expected:    nil,
		},
	}

	db, err := dbtest.New()
	if err != nil {
		t.Fatalf("Failed to create the test container: %s", err.Error())
	}
	defer db.Teardown()

	store := New(context.TODO(), db.URI)

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, db.ApplyFixtures(tc.fixtures...))
			defer db.TeardownFixtures()

			err := store.UserSetConfirmed(context.TODO(), tc.id, tc.confirmed)
			assert.Equal(t, tc.expected, err)
		})
	}
}
