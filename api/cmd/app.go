package main

import (
	"context"
	"os"
	"reflect"
	"time"

	"github.com/heiytor/invenda/api/pkg/env"
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

	if err := env.Load(); err != nil {
		log.Panic().
			Err(err).
			Msg("Unable to parse the environment variables")
	}

	v := reflect.ValueOf(env.E())
	for i := 0; i < v.NumField(); i++ {
		log.Info().
			Str("key", v.Type().Field(i).Name).
			Interface("value", v.Field(i).Interface()).
			Msg("Environment variable loaded")
	}

	client, preferredDb, err := store.Connect(ctx, env.E().MongoURI)
	if err != nil {
		log.Panic().
			Err(err).
			Msg("Unable to connect to MongoDB")
	}

	log.Info().Msg("Connected to MongoDB")

	store, err := store.New(ctx, client, preferredDb)
	if err != nil {
		log.Panic().
			Err(err).
			Msg("Unable to create the store")
	}

	service := service.New(store)
	routes := route.New(service)

	// Configure logger
	logger := lecho.From(log.Logger)
	routes.E.Logger = logger
	routes.E.Use(middleware.Logger(logger))

	if err := routes.E.Start(":8080"); err != nil {
		log.Panic().
			Err(err).
			Msg("Echo panicked.")
	}
}
