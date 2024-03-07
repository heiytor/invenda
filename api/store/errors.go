package store

import (
	"errors"
	"io"

	"go.mongodb.org/mongo-driver/mongo"
)

// fromMongoError converte um erro do mongo em um erro conhecido e gerenciado pela aplicação.
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

var (
	ErrUnexpected = errors.New("Unexpected Error")
	ErrNotFound   = errors.New("Document not found")
)
