package main

import (
	"os"
	"time"

	"github.com/heiytor/invenda/api/route"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v3"
)

func main() {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	log.Logger = logger

	e := route.New(lecho.From(logger))

	e.Logger.Fatal(e.Start(":3333"))
}
