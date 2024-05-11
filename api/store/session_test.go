package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/heiytor/invenda/api/pkg/query"
	"github.com/heiytor/invenda/api/store"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestSessionGet(t *testing.T) {
	type Actual struct {
		session *models.Session
		err     error
	}

	cases := []struct {
		description string
		id          string
		opts        []store.GetSessionOption
		fixtures    []fixture
		expected    Actual
	}{
		{
			description: "fails when session is not found",
			id:          "00000000000000000000000000",
			opts:        []store.GetSessionOption{},
			fixtures:    []fixture{},
			expected: Actual{
				session: nil,
				err:     store.ErrNotFound,
			},
		},
		{
			description: "succeeds to find a session",
			id:          "ss_01HX6FABC3SRPVK6VSM48DWMMQ",
			opts:        []store.GetSessionOption{},
			fixtures:    []fixture{fixtureSession},
			expected: Actual{
				session: &models.Session{
					ID:        "ss_01HX6FABC3SRPVK6VSM48DWMMQ",
					UserID:    "usr_01HNGJ2BTGQAHAZ1XNYZQPG719",
					StartedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					EndedAt:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					State:     models.SessionStateActive,
					SourceIP:  "127.0.0.1",
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

			session, err := s.Session.Get(ctx, tc.id, tc.opts...)
			require.Equal(t, tc.expected, Actual{session, err})
		})
	}
}

func TestSessionList(t *testing.T) {
	type Actual struct {
		session []models.Session
		count   int64
		err     error
	}

	cases := []struct {
		description string
		userID      string
		query       *query.Query
		opts        []store.GetSessionOption
		fixtures    []fixture
		expected    Actual
	}{
		{
			description: "succeeds when session is not found",
			userID:      "usr_00000000000000000000000000",
			query:       &query.Query{Sorter: query.Sorter{By: "started_at", Order: query.OrderAsc}},
			opts:        []store.GetSessionOption{},
			fixtures:    []fixture{},
			expected: Actual{
				session: []models.Session{},
				count:   0,
				err:     nil,
			},
		},
		{
			description: "succeeds to list the sessions with order asc",
			userID:      "usr_01HNGJ2BTGQAHAZ1XNYZQPG719",
			query:       &query.Query{Sorter: query.Sorter{By: "started_at", Order: query.OrderAsc}},
			opts:        []store.GetSessionOption{},
			fixtures:    []fixture{fixtureSession},
			expected: Actual{
				session: []models.Session{
					{
						ID:        "ss_01HX6FABC3SRPVK6VSM48DWMMQ",
						UserID:    "usr_01HNGJ2BTGQAHAZ1XNYZQPG719",
						StartedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						EndedAt:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						State:     models.SessionStateActive,
						SourceIP:  "127.0.0.1",
					},
					{
						ID:        "ss_01HX6FP9QED8TW64DGZHMYDWVF",
						UserID:    "usr_01HNGJ2BTGQAHAZ1XNYZQPG719",
						StartedAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
						EndedAt:   time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
						State:     models.SessionStateInactive,
						SourceIP:  "127.0.0.1",
					},
				},
				count: 2,
				err:   nil,
			},
		},
		{
			description: "succeeds to list the sessions with order desc",
			userID:      "usr_01HNGJ2BTGQAHAZ1XNYZQPG719",
			query:       &query.Query{Sorter: query.Sorter{By: "started_at", Order: query.OrderDesc}},
			opts:        []store.GetSessionOption{},
			fixtures:    []fixture{fixtureSession},
			expected: Actual{
				session: []models.Session{
					{
						ID:        "ss_01HX6FP9QED8TW64DGZHMYDWVF",
						UserID:    "usr_01HNGJ2BTGQAHAZ1XNYZQPG719",
						StartedAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
						EndedAt:   time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
						State:     models.SessionStateInactive,
						SourceIP:  "127.0.0.1",
					},
					{
						ID:        "ss_01HX6FABC3SRPVK6VSM48DWMMQ",
						UserID:    "usr_01HNGJ2BTGQAHAZ1XNYZQPG719",
						StartedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						EndedAt:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						State:     models.SessionStateActive,
						SourceIP:  "127.0.0.1",
					},
				},
				count: 2,
				err:   nil,
			},
		},
		{
			description: "succeeds to list the sessions with pagination",
			userID:      "usr_01HNGJ2BTGQAHAZ1XNYZQPG719",
			query:       &query.Query{Paginator: query.Paginator{Page: 1, Size: 1}, Sorter: query.Sorter{By: "started_at", Order: query.OrderAsc}},
			opts:        []store.GetSessionOption{},
			fixtures:    []fixture{fixtureSession},
			expected: Actual{
				session: []models.Session{
					{
						ID:        "ss_01HX6FABC3SRPVK6VSM48DWMMQ",
						UserID:    "usr_01HNGJ2BTGQAHAZ1XNYZQPG719",
						StartedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						EndedAt:   time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
						State:     models.SessionStateActive,
						SourceIP:  "127.0.0.1",
					},
				},
				count: 2,
				err:   nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			srv.apply(tc.fixtures...)
			defer srv.reset()

			ctx := context.Background()

			sessions, count, err := s.Session.List(ctx, tc.userID, tc.query, tc.opts...)
			require.Equal(t, tc.expected, Actual{sessions, count, err})
		})
	}
}

func TestSessionCreate(t *testing.T) {
	type Actual struct {
		err error
	}

	cases := []struct {
		description string
		session     *models.Session
		expected    Actual
	}{
		{
			description: "succeeds to create a session",
			session:     &models.Session{},
			expected:    Actual{err: nil},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			defer srv.reset()
			ctx := context.Background()

			id, err := s.Session.Create(ctx, tc.session)
			require.Equal(t, tc.expected, Actual{err})
			require.NotEmpty(t, id)

			// TODO:
		})
	}
}

func TestSessionUpdate(t *testing.T) {
	type Actual struct {
		err error
	}

	cases := []struct {
		description string
		id          string
		changes     *models.SessionChanges
		fixtures    []fixture
		expected    Actual
	}{
		{
			description: "fails when session is not found",
			id:          "ss_00000000000000000000000000",
			changes:     &models.SessionChanges{},
			fixtures:    []fixture{},
			expected:    Actual{err: store.ErrNotFound},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			srv.apply(tc.fixtures...)
			defer srv.reset()

			ctx := context.Background()

			if err := s.Session.Update(ctx, tc.id, tc.changes); err != nil {
				require.Equal(t, tc.expected, Actual{err})
				return
			}

			// TODO:
		})
	}
}

func TestSessionDelete(t *testing.T) {
	type Actual struct {
		err error
	}

	cases := []struct {
		description string
		id          string
		fixtures    []fixture
		expected    Actual
	}{
		{
			description: "fails when session is not found",
			id:          "ss_00000000000000000000000000",
			fixtures:    []fixture{},
			expected:    Actual{err: store.ErrNotFound},
		},
		{
			description: "succeeds to delete a session",
			fixtures:    []fixture{fixtureSession},
			id:          "ss_01HX6FABC3SRPVK6VSM48DWMMQ",
			expected:    Actual{err: nil},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			srv.apply(tc.fixtures...)
			defer srv.reset()

			ctx := context.Background()

			if err := s.Session.Delete(ctx, tc.id); err != nil {
				require.Equal(t, tc.expected, Actual{err})
				return
			}

			session := new(models.Session)
			require.Error(t, db.Collection("session").FindOne(ctx, bson.M{"_id": tc.id}).Decode(session))
		})
	}
}
