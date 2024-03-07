package dbtest

import (
	"context"
	"path/filepath"
	"runtime"

	"github.com/shellhub-io/mongotest"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

type DB struct {
	Container *mongodb.MongoDBContainer
	URI       string
}

func New() (*DB, error) {
	ctx := context.TODO()

	container, err := mongodb.RunContainer(ctx, testcontainers.WithImage("mongo:7.0.5"))
	if err != nil {
		return nil, err
	}

	uri, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	// Habilita o uso de fixtures
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to retrieve the fixtures path at runtime")
	}

	mongotest.Configure(mongotest.Config{
		URL:            uri,
		Database:       "test",
		FixtureRootDir: filepath.Join(filepath.Dir(file), "fixtures"),
		FixtureFormat:  mongotest.FixtureFormatJSON,
		PreInsertFuncs: []mongotest.PreInsertFunc{
			mongotest.SimpleConvertTime("user", "created_at"),
			mongotest.SimpleConvertTime("user", "updated_at"),
		},
	})

	return &DB{container, uri + "/test"}, nil
}

func (db *DB) Teardown() {
	if err := db.Container.Terminate(context.TODO()); err != nil {
		panic(err)
	}
}
