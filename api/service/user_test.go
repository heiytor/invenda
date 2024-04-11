package service_test

import (
	"context"
	goerrors "errors"
	"net/http"
	"testing"
	"time"

	"github.com/heiytor/invenda/api/pkg/errors"
	"github.com/heiytor/invenda/api/pkg/hash"
	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/heiytor/invenda/api/pkg/requests"
	"github.com/heiytor/invenda/api/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	type Actual struct {
		insertedID string
		err        error
	}

	cases := []struct {
		description string
		req         *requests.CreateUser
		mocks       func(context.Context, *mocks)
		expected    Actual
	}{
		{
			description: "fails when req has conflicted fiels",
			req: &requests.CreateUser{
				Name:     "John Doe",
				Email:    "john.doe@test.com",
				Password: "secret",
			},
			mocks: func(ctx context.Context, m *mocks) {
				m.user.
					On("Conflicts", ctx, &models.User{Email: "john.doe@test.com"}).
					Return([]string{"email"}, nil).
					Once()
			},
			expected: Actual{
				insertedID: "",
				err: &errors.Error{
					Code:  http.StatusConflict,
					Layer: errors.LayerService,
					Msg:   errors.MsgConflict,
					Attrs: map[string]interface{}{
						"entity":    "user",
						"conflicts": []string{"email"},
					},
				},
			},
		},
		{
			description: "fails when unable to create the user",
			req: &requests.CreateUser{
				Name:     "John Doe",
				Email:    "john.doe@test.com",
				Password: "secret",
			},
			mocks: func(ctx context.Context, m *mocks) {
				m.user.
					On("Conflicts", ctx, &models.User{Email: "john.doe@test.com"}).
					Return([]string{}, nil).
					Once()
				m.hash.
					On("New", "secret", hash.DefaultOptions).
					Return("$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0").
					Once()
				m.user.
					On("Create", ctx, &models.User{
						Name:     "John Doe",
						Email:    "john.doe@test.com",
						Password: "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0",
					}).
					Return("", goerrors.New("error")).
					Once()
			},
			expected: Actual{
				insertedID: "",
				err:        goerrors.New("error"),
			},
		},
		{
			description: "succeeds when able to create the user",
			req: &requests.CreateUser{
				Name:     "John Doe",
				Email:    "john.doe@test.com",
				Password: "secret",
			},
			mocks: func(ctx context.Context, m *mocks) {
				m.user.
					On("Conflicts", ctx, &models.User{Email: "john.doe@test.com"}).
					Return([]string{}, nil).
					Once()
				m.hash.
					On("New", "secret", hash.DefaultOptions).
					Return("$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0").
					Once()
				m.user.
					On("Create", ctx, &models.User{
						Name:     "John Doe",
						Email:    "john.doe@test.com",
						Password: "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0",
					}).
					Return("usr_01HV3JEM0X9XQN6HQ7685SM4Z7", nil).
					Once()
			},
			expected: Actual{
				insertedID: "usr_01HV3JEM0X9XQN6HQ7685SM4Z7",
				err:        nil,
			},
		},
	}

	m := getMocks(t)
	s := service.New(m.store)
	defer m.RequireExpectations()

	for _, tt := range cases {
		tc := tt
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()
			tc.mocks(ctx, m)

			insertedID, err := s.CreateUser(ctx, tc.req)
			require.Equal(t, tc.expected, Actual{insertedID, err})
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type Actual struct {
		user *models.User
		err  error
	}

	cases := []struct {
		description string
		id          string
		req         *requests.UpdateUser
		mocks       func(context.Context, *mocks)
		expected    Actual
	}{
		{
			description: "fails when req has conflicted fiels",
			id:          "usr_01HV3JEM0X9XQN6HQ7685SM4Z7",
			req: &requests.UpdateUser{
				Name:     "John Doe",
				Email:    "john.doe@test.com",
				Password: "secret",
			},
			mocks: func(ctx context.Context, m *mocks) {
				m.user.
					On("Conflicts", ctx, &models.User{Email: "john.doe@test.com"}).
					Return([]string{"email"}, nil).
					Once()
			},
			expected: Actual{
				user: nil,
				err: &errors.Error{
					Code:  http.StatusConflict,
					Layer: errors.LayerService,
					Msg:   errors.MsgConflict,
					Attrs: map[string]interface{}{
						"entity":    "user",
						"conflicts": []string{"email"},
					},
				},
			},
		},
		{
			description: "fails when unable to update the user",
			id:          "usr_01HV3JEM0X9XQN6HQ7685SM4Z7",
			req: &requests.UpdateUser{
				Name:     "John Doe",
				Email:    "john.doe@test.com",
				Password: "secret",
			},
			mocks: func(ctx context.Context, m *mocks) {
				m.user.
					On("Conflicts", ctx, &models.User{Email: "john.doe@test.com"}).
					Return([]string{}, nil).
					Once()
				m.hash.
					On("New", "secret", hash.DefaultOptions).
					Return("$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0").
					Once()
				m.user.
					On("Update", ctx, "usr_01HV3JEM0X9XQN6HQ7685SM4Z7", &models.UserChanges{
						Name:     "John Doe",
						Email:    "john.doe@test.com",
						Password: "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0",
					}).
					Return(goerrors.New("error")).
					Once()
			},
			expected: Actual{
				user: nil,
				err:  goerrors.New("error"),
			},
		},
		{
			description: "fails when unable to retrieve the updated user",
			id:          "usr_01HV3JEM0X9XQN6HQ7685SM4Z7",
			req: &requests.UpdateUser{
				Name:     "John Doe",
				Email:    "john.doe@test.com",
				Password: "secret",
			},
			mocks: func(ctx context.Context, m *mocks) {
				m.user.
					On("Conflicts", ctx, &models.User{Email: "john.doe@test.com"}).
					Return([]string{}, nil).
					Once()
				m.hash.
					On("New", "secret", hash.DefaultOptions).
					Return("$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0").
					Once()
				m.user.
					On("Update", ctx, "usr_01HV3JEM0X9XQN6HQ7685SM4Z7", &models.UserChanges{
						Name:     "John Doe",
						Email:    "john.doe@test.com",
						Password: "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0",
					}).
					Return(nil).
					Once()
				m.user.
					On("GetByID", ctx, "usr_01HV3JEM0X9XQN6HQ7685SM4Z7", mock.Anything).
					Return(nil, goerrors.New("error")).
					Once()
			},
			expected: Actual{
				user: nil,
				err:  goerrors.New("error"),
			},
		},
		{
			description: "succeeds to update and retrieve the user",
			id:          "usr_01HV3JEM0X9XQN6HQ7685SM4Z7",
			req: &requests.UpdateUser{
				Name:     "John Doe",
				Email:    "john.doe@test.com",
				Password: "secret",
			},
			mocks: func(ctx context.Context, m *mocks) {
				m.user.
					On("Conflicts", ctx, &models.User{Email: "john.doe@test.com"}).
					Return([]string{}, nil).
					Once()
				m.hash.
					On("New", "secret", hash.DefaultOptions).
					Return("$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0").
					Once()
				m.user.
					On("Update", ctx, "usr_01HV3JEM0X9XQN6HQ7685SM4Z7", &models.UserChanges{
						Name:     "John Doe",
						Email:    "john.doe@test.com",
						Password: "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0",
					}).
					Return(nil).
					Once()
				m.user.
					On("GetByID", ctx, "usr_01HV3JEM0X9XQN6HQ7685SM4Z7", mock.Anything).
					Return(
						&models.User{
							ID:        "usr_01HV3JEM0X9XQN6HQ7685SM4Z7",
							CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
							UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
							Name:      "John Doe",
							Email:     "john.doe@test.com",
							Password:  "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0",
						},
						nil,
					).
					Once()
			},
			expected: Actual{
				user: &models.User{
					ID:        "usr_01HV3JEM0X9XQN6HQ7685SM4Z7",
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					Name:      "John Doe",
					Email:     "john.doe@test.com",
					Password:  "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0",
				},
				err: nil,
			},
		},
	}

	m := getMocks(t)
	s := service.New(m.store)
	defer m.RequireExpectations()

	for _, tt := range cases {
		tc := tt
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()
			tc.mocks(ctx, m)

			usr, err := s.UpdateUser(ctx, tc.id, tc.req)
			require.Equal(t, tc.expected, Actual{usr, err})
		})
	}
}

func TestAuthUser(t *testing.T) {
	type Actual struct {
		u *models.User
		t *models.AccessToken
		e error
	}

	m := getMocks(t)
	defer m.RequireExpectations()

	cases := []struct {
		description string
		req         *requests.AuthUser
		mocks       func(context.Context, *mocks)
		expected    Actual
	}{
		{
			description: "fails when unable to retrieve the user",
			req: &requests.AuthUser{
				Identifier: "john.doe@test.com",
				Password:   "secret",
			},
			mocks: func(ctx context.Context, m *mocks) {
				m.user.
					On("GetByEmail", ctx, "john.doe@test.com").
					Return(nil, goerrors.New("error")).
					Once()
			},
			expected: Actual{
				u: nil,
				t: nil,
				e: goerrors.New("error"),
			},
		},
		{
			description: "fails when req.Password does not match with usr.Password",
			req: &requests.AuthUser{
				Identifier: "john.doe@test.com",
				Password:   "secret",
			},
			mocks: func(ctx context.Context, m *mocks) {
				m.user.
					On("GetByEmail", ctx, "john.doe@test.com").
					Return(
						&models.User{
							ID:        "usr_01HV3JEM0X9XQN6HQ7685SM4Z7",
							CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
							UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
							Name:      "John Doe",
							Email:     "john.doe@test.com",
							Password:  "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0",
						},
						nil,
					).
					Once()
				m.hash.
					On("Compare", "secret", "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0").
					Return(false).
					Once()
			},
			expected: Actual{
				u: nil,
				t: nil,
				e: &errors.Error{
					Code:  500,
					Msg:   "foo",
					Attrs: map[string]interface{}{},
				},
			},
		},
		{
			description: "fails when req.Password does not match with usr.Password",
			req: &requests.AuthUser{
				Identifier: "john.doe@test.com",
				Password:   "secret",
			},
			mocks: func(ctx context.Context, m *mocks) {
				m.user.
					On("GetByEmail", ctx, "john.doe@test.com").
					Return(
						&models.User{
							ID:        "usr_01HV3JEM0X9XQN6HQ7685SM4Z7",
							CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
							UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
							Name:      "John Doe",
							Email:     "john.doe@test.com",
							Password:  "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0",
						},
						nil,
					).
					Once()
				m.hash.
					On("Compare", "secret", "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0").
					Return(true).
					Once()
				m.clock.
					On("Now").
					Return(m.now).
					Once()
				m.user.
					On("CreateRefreshToken", ctx, &models.AccessToken{
						UserID:    "usr_01HV3JEM0X9XQN6HQ7685SM4Z7",
						Status:    "active",
						ExpiresIn: m.now.AddDate(0, 1, 0).Unix(),
					}).
					Return("", goerrors.New("error")).
					Once()
			},
			expected: Actual{
				u: nil,
				t: nil,
				e: goerrors.New("error"),
			},
		},
		{
			description: "succeeds to create the user's token",
			req: &requests.AuthUser{
				Identifier: "john.doe@test.com",
				Password:   "secret",
			},
			mocks: func(ctx context.Context, m *mocks) {
				m.user.
					On("GetByEmail", ctx, "john.doe@test.com").
					Return(
						&models.User{
							ID:        "usr_01HV3JEM0X9XQN6HQ7685SM4Z7",
							CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
							UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
							Name:      "John Doe",
							Email:     "john.doe@test.com",
							Password:  "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0",
						},
						nil,
					).
					Once()
				m.hash.
					On("Compare", "secret", "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0").
					Return(true).
					Once()
				m.clock.
					On("Now").
					Return(m.now).
					Once()
				m.user.
					On("CreateRefreshToken", ctx, &models.AccessToken{
						UserID:    "usr_01HV3JEM0X9XQN6HQ7685SM4Z7",
						Status:    "active",
						ExpiresIn: m.now.AddDate(0, 1, 0).Unix(),
					}).
					Return("usr-tkn_01HV3N01M4NE85MS4Y8GHMDG5B", nil).
					Once()
			},
			expected: Actual{
				u: &models.User{
					ID:        "usr_01HV3JEM0X9XQN6HQ7685SM4Z7",
					CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					Name:      "John Doe",
					Email:     "john.doe@test.com",
				},
				t: &models.AccessToken{
					UserID:    "usr_01HV3JEM0X9XQN6HQ7685SM4Z7",
					Status:    "active",
					ExpiresIn: m.now.AddDate(0, 1, 0).Unix(),
				},
				e: nil,
			},
		},
	}

	s := service.New(m.store)

	for _, tt := range cases {
		tc := tt
		t.Run(tc.description, func(token *testing.T) {
			ctx := context.Background()
			tc.mocks(ctx, m)

			usr, tkn, err := s.AuthUser(ctx, tc.req)
			require.Equal(token, tc.expected, Actual{usr, tkn, err})
		})
	}
}
