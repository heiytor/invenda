package main

import (
	"context"
	"os"
	"time"

	"github.com/heiytor/invenda/api/route"
	"github.com/heiytor/invenda/api/service"
	"github.com/heiytor/invenda/api/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v3"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	ctx := context.Background()

	store := store.New(ctx)
	service := service.New(store)
	routes := route.New(service, lecho.From(log.Logger))

	if err := routes.E.Start(":3333"); err != nil {
		log.Panic().Err(err).Msg("Echo panicked.")
	}
}
