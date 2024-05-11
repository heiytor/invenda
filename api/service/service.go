package service

import (
	"github.com/heiytor/invenda/api/pkg/cache"
	"github.com/heiytor/invenda/api/store"
)

type service struct {
	store *store.Store
	cache cache.Cache
}

type Service interface {
	User
	Namespace
	Session
}

func New(store *store.Store, cache cache.Cache) Service {
	return &service{store: store, cache: cache}
}
