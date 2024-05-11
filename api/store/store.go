package store

import (
	"context"

	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Store struct {
	client *mongodb.Client
	db     *mongodb.Database

	// User handle all user-related operations.
	User      User
	Namespace Namespace
	Session   Session
}

func Connect(ctx context.Context, uri string) (*mongodb.Client, string, error) {
	connstr, err := connstring.ParseAndValidate(uri)
	if err != nil {
		return nil, "", err
	}

	client, err := mongodb.Connect(ctx, options.Client().ApplyURI(connstr.Original))
	if err != nil {
		return nil, "", err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, "", err
	}

	return client, connstr.Database, nil
}

func New(ctx context.Context, client *mongodb.Client, db string) (*Store, error) {
	store := &Store{
		client: client,
		db:     client.Database(db),
	}

	store.User = &user{c: store.db.Collection("user")}
	store.Namespace = &namespace{c: store.db.Collection("namespace")}
	store.Session = &session{c: store.db.Collection("session")}

	return store, nil
}
