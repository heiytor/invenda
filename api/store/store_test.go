package store_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/heiytor/invenda/api/store"
	"github.com/shellhub-io/mongotest"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	Container *mongodb.MongoDBContainer
	URI       string
}

// TODO
var srv *Server

// TODO
var s *store.Store

// TODO
var db *mongo.Database

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := mongodb.RunContainer(ctx, testcontainers.WithImage("mongo:7.0.5"))
	if err != nil {
		panic(err)
	}
	defer container.Terminate(ctx)

	uri, err := container.ConnectionString(ctx)
	if err != nil {
		panic(err)
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

	srv = &Server{container, uri + "/test"}

	c, preferredDb, err := store.Connect(ctx, srv.URI)
	if err != nil {
		panic(err)
	}

	db = c.Database(preferredDb)
	s, err = store.New(ctx, c, preferredDb)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

type fixture string

const (
	fixtureUser fixture = "user"
)

func (*Server) apply(fixtures ...fixture) error {
	var str []string
	for _, f := range fixtures {
		str = append(str, string(f))
	}

	return mongotest.UseFixture(str...)
}

func (*Server) reset(fixtures ...string) error {
	return mongotest.DropDatabase()
}
