package store

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type store struct {
	client *mongo.Client
	db     *mongo.Database
	users  *mongo.Collection // users é um wrapper para a coleção "user"
}

type Store interface {
	User
}

func New(ctx context.Context, conn string) Store {
	client, db, err := connectDB(ctx, conn)
	if err != nil {
		panic(err)
	}

	store := &store{client: client, db: db}
	store.users = db.Collection("user")

	return store
}

func connectDB(ctx context.Context, conn string) (*mongo.Client, *mongo.Database, error) {
	connStr, err := connstring.ParseAndValidate(conn)
	if err != nil {
		return nil, nil, err
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr.Original))
	if err != nil {
		return nil, nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}

	return client, client.Database(connStr.Database), nil
}
