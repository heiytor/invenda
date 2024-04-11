package store

import (
	"errors"
	"io"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUnexpected = errors.New("unexpected Error")
	ErrNotFound   = errors.New("document not found")
)

func fromMongoError(err error) error {
	switch {
	case err == mongo.ErrNoDocuments, err == io.EOF:
		return ErrNotFound
	default:
		if err == nil {
			return nil
		}

		return errors.Join(ErrUnexpected, err)
	}
}
