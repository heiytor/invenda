package service

import (
	"github.com/heiytor/invenda/api/store"
)

type service struct {
	store store.Store
}

type Service interface{}

func New(store store.Store) Service {
	return &service{store: store}
}
