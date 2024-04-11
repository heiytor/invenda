package main

import (
	"context"
	"os"
	"time"

	"github.com/heiytor/invenda/api/pkg/secretkeys"
	"github.com/heiytor/invenda/api/route"
	"github.com/heiytor/invenda/api/route/pkg/middleware"
	"github.com/heiytor/invenda/api/service"
	"github.com/heiytor/invenda/api/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v3"
)

func main() {
	ctx := context.Background()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	if err := secretkeys.Load(); err != nil {
		log.Panic().
			Err(err).
			Msg("Unable to parse the secret keys")
	}

	client, preferredDb, err := store.Connect(ctx, "mongodb://mongo:27017/main")
	if err != nil {
		log.Panic().
			Err(err).
			Msg("Unable to connect to MongoDB")
	}

	store, err := store.New(ctx, client, preferredDb)
	if err != nil {
		log.Panic().
			Err(err).
			Msg("Unable to create the store")
	}

	log.Info().Msg("Connected to MongoDB")

	routes := route.New(service.New(store))

	// Configure logger
	logger := lecho.From(log.Logger)
	routes.E.Logger = logger
	routes.E.Use(middleware.Logger(logger))

	if err := routes.E.Start(":3333"); err != nil {
		log.Panic().
			Err(err).
			Msg("Echo panicked.")
	}
}
