package service_test

import (
	"os"
	"testing"
	"time"

	"github.com/heiytor/invenda/api/pkg/clock"
	clockmocks "github.com/heiytor/invenda/api/pkg/clock/mocks"
	"github.com/heiytor/invenda/api/pkg/hash"
	hashmocks "github.com/heiytor/invenda/api/pkg/hash/mocks"
	"github.com/heiytor/invenda/api/store"
	storemocks "github.com/heiytor/invenda/api/store/mocks"
	"github.com/stretchr/testify/require"
)

type mocks struct {
	t   *testing.T
	now time.Time

	store *store.Store
	user  *storemocks.User
	hash  *hashmocks.Hash
	clock *clockmocks.Clock
}

func getMocks(t *testing.T) *mocks {
	userMock := &storemocks.User{}

	hashMock := &hashmocks.Hash{}
	hash.Backend = hashMock

	clockMock := &clockmocks.Clock{}
	clock.Backend = clockMock

	s := &store.Store{}
	s.User = userMock

	return &mocks{
		t:     t,
		store: s,
		now:   time.Now().UTC(),
		user:  userMock,
		hash:  hashMock,
		clock: clockMock,
	}
}

func (m *mocks) RequireExpectations() {
	require.True(m.t, m.clock.AssertExpectations(m.t))
	require.True(m.t, m.user.AssertExpectations(m.t))
	require.True(m.t, m.hash.AssertExpectations(m.t))
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
