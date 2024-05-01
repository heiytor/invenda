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

func TestNamespaceGet(t *testing.T) {
	type Actual struct {
		namespace *models.Namespace
		err       error
	}

	cases := []struct {
		description string
		id          string
		opts        []store.GetNamespaceOption
		fixtures    []fixture
		expected    Actual
	}{
		{
			description: "fails when namespace is not found",
			id:          "00000000000000000000000000",
			opts:        []store.GetNamespaceOption{},
			fixtures:    []fixture{},
			expected: Actual{
				namespace: nil,
				err:       store.ErrNotFound,
			},
		},
		{
			description: "succeeds to find a namespace",
			id:          "ns_01HV7FKH5SRB0TGWM7MQ15PYKN",
			opts:        []store.GetNamespaceOption{},
			fixtures:    []fixture{fixtureNamespace},
			expected: Actual{
				namespace: &models.Namespace{
					ID:        "ns_01HV7FKH5SRB0TGWM7MQ15PYKN",
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					Name:      "Atacarejo",
					Members: []models.Member{
						{
							ID:          "usr_01HNGJ2BTGQAHAZ1XNYZQPG719",
							AddedAt:     time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
							Permissions: []string{"namespace:write"},
						},
					},
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

			namespace, err := s.Namespace.Get(ctx, tc.id, tc.opts...)
			require.Equal(t, tc.expected, Actual{namespace, err})
		})
	}
}

func TestNamespaceCreate(t *testing.T) {
	type Actual struct {
		err error
	}

	cases := []struct {
		description string
		namespace   *models.Namespace
		expected    Actual
	}{
		{
			description: "succeeds to create a namespace",
			namespace: &models.Namespace{
				Name: "Atarecejo",
				Members: []models.Member{
					{
						ID:          "01HNGJ2BTGQAHAZ1XNYZQPG719",
						Permissions: []string{"namespace:write"},
					},
				},
			},
			expected: Actual{err: nil},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			defer srv.reset()
			ctx := context.Background()

			id, err := s.Namespace.Create(ctx, tc.namespace)
			require.Equal(t, tc.expected, Actual{err})
			require.NotEmpty(t, id)

			namespace := new(models.Namespace)
			require.NoError(t, db.Collection("namespace").FindOne(ctx, bson.M{"_id": id}).Decode(namespace))
			require.Equal(t, tc.namespace.Name, namespace.Name)
			require.Equal(t, len(tc.namespace.Members), 1)
			require.Equal(t, tc.namespace.Members[0].ID, namespace.Members[0].ID)
			require.Equal(t, tc.namespace.Members[0].Permissions[0], namespace.Members[0].Permissions[0])
		})
	}
}

func TestNamespaceUpdate(t *testing.T) {
	type Actual struct {
		err error
	}

	cases := []struct {
		description string
		id          string
		changes     *models.NamespaceChanges
		fixtures    []fixture
		expected    Actual
	}{
		{
			description: "fails when namespace is not found",
			id:          "ns_00000000000000000000000000",
			changes:     &models.NamespaceChanges{},
			fixtures:    []fixture{},
			expected:    Actual{err: store.ErrNotFound},
		},
		{
			description: "succeeds to update a namespace",
			id:          "ns_01HV7FKH5SRB0TGWM7MQ15PYKN",
			changes:     &models.NamespaceChanges{Name: "New Name"},
			fixtures:    []fixture{fixtureNamespace},
			expected:    Actual{err: nil},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			srv.apply(tc.fixtures...)
			defer srv.reset()

			ctx := context.Background()

			if err := s.Namespace.Update(ctx, tc.id, tc.changes); err != nil {
				require.Equal(t, tc.expected, Actual{err})
				return
			}

			namespace := new(models.Namespace)
			require.NoError(t, db.Collection("namespace").FindOne(ctx, bson.M{"_id": tc.id}).Decode(namespace))
			require.Equal(t, tc.changes.Name, namespace.Name)
			require.WithinDuration(t, tc.changes.UpdatedAt, namespace.UpdatedAt, time.Millisecond)
		})
	}
}

func TestNamespaceDelete(t *testing.T) {
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
			description: "fails when namespace is not found",
			id:          "ns_00000000000000000000000000",
			fixtures:    []fixture{},
			expected:    Actual{err: store.ErrNotFound},
		},
		{
			description: "succeeds to delete a namespace",
			fixtures:    []fixture{fixtureNamespace},
			id:          "ns_01HV7FKH5SRB0TGWM7MQ15PYKN",
			expected:    Actual{err: nil},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			srv.apply(tc.fixtures...)
			defer srv.reset()

			ctx := context.Background()

			if err := s.Namespace.Delete(ctx, tc.id); err != nil {
				require.Equal(t, tc.expected, Actual{err})
				return
			}

			namespace := new(models.Namespace)
			require.Error(t, db.Collection("namespace").FindOne(ctx, bson.M{"_id": tc.id}).Decode(namespace))
		})
	}
}

func TestNamespaceUpsertMember(t *testing.T) {
	type Actual struct {
		err error
	}

	cases := []struct {
		description string
		id          string
		member      *models.Member
		exists      bool
		fixtures    []fixture
		expected    Actual
	}{
		{
			description: "fails when namespace is not found",
			id:          "ns_00000000000000000000000000",
			member: &models.Member{
				ID:          "usr_01HVD0VWF0G6Z7TMCP1CPST4FN",
				Permissions: []string{"namespace:read", "namepace:write"},
			},
			exists:   false,
			fixtures: []fixture{},
			expected: Actual{err: store.ErrNotFound},
		},
		{
			description: "succeeds to add a member",
			id:          "ns_01HV7FKH5SRB0TGWM7MQ15PYKN",
			member: &models.Member{
				ID:          "usr_01HVD0VWF0G6Z7TMCP1CPST4FN",
				Permissions: []string{"namespace:read", "namepace:write"},
			},
			exists:   false,
			fixtures: []fixture{fixtureNamespace},
			expected: Actual{err: nil},
		},
		{
			description: "succeeds to add a member",
			id:          "ns_01HV7FKH5SRB0TGWM7MQ15PYKN",
			member: &models.Member{
				ID:          "usr_01HNGJ2BTGQAHAZ1XNYZQPG719",
				AddedAt:     time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				Permissions: []string{"namespace:read", "namepace:write"},
			},
			exists:   true,
			fixtures: []fixture{fixtureNamespace},
			expected: Actual{err: nil},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			srv.apply(tc.fixtures...)
			defer srv.reset()

			ctx := context.Background()

			if err := s.Namespace.UpsertMember(ctx, tc.id, tc.member, tc.exists); err != nil {
				require.Equal(t, tc.expected, Actual{err})
				return
			}

			namespace := new(models.Namespace)
			require.NoError(t, db.Collection("namespace").FindOne(ctx, bson.M{"_id": tc.id}).Decode(namespace))

			for _, m := range namespace.Members {
				if m.ID == tc.member.ID {
					require.Equal(t, tc.member.Permissions, namespace.Members[1].Permissions)
				}
			}
		})
	}
}

func TestNamespaceRemoveMember(t *testing.T) {
	type Actual struct {
		err error
	}

	cases := []struct {
		description string
		id          string
		memberID    string
		fixtures    []fixture
		expected    Actual
	}{
		{
			description: "fails when namespace is not found",
			id:          "ns_00000000000000000000000000",
			memberID:    "usr_01HNGJ2BTGQAHAZ1XNYZQPG719",
			fixtures:    []fixture{},
			expected:    Actual{err: store.ErrNotFound},
		},
		{
			description: "succeeds to add a member",
			id:          "ns_01HV7FKH5SRB0TGWM7MQ15PYKN",
			memberID:    "usr_01HNGJ2BTGQAHAZ1XNYZQPG719",
			fixtures:    []fixture{fixtureNamespace},
			expected:    Actual{err: nil},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			srv.apply(tc.fixtures...)
			defer srv.reset()

			ctx := context.Background()

			if err := s.Namespace.RemoveMember(ctx, tc.id, tc.memberID); err != nil {
				require.Equal(t, tc.expected, Actual{err})
				return
			}

			namespace := new(models.Namespace)
			require.NoError(t, db.Collection("namespace").FindOne(ctx, bson.M{"_id": tc.id}).Decode(namespace))
			require.Equal(t, 0, len(namespace.Members))
		})
	}
}
