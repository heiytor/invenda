package env

import (
	"context"

	"github.com/heiytor/invenda/api/pkg/validator"
	"github.com/sethvargo/go-envconfig"
)

type spec struct {
	// Version specifies the current Invenda version.
	Version string `env:"INVENDA_VERSION" validate:"required"`
	// Environment specifies whether Invenda is running in a development or production environment.
	Environment string `env:"INVENDA_ENVIRONMENT" validate:"required|in:development,production"`
	// MongoURI stores the connection URI for MongoDB.
	MongoURI string `env:"INVENDA_MONGO_URI" validate:"required"`
}

var s = new(spec)

// Load initializes the environment variables based on the provided configuration specification.
// This function must be called before using any other functions in this package.
func Load() error {
	if err := envconfig.Process(context.Background(), s); err != nil {
		return err
	}

	if err := validator.New().Validate(s); err != nil {
		return err
	}

	return nil
}

// E returns a copy of `spec` instance that contains the loaded environment variables.
// The function should be called after calling the [Load] function to ensure the variables
// have been initialized and validated.
func E() spec {
	return *s
}
