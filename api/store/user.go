package store

import (
	"context"

	"github.com/heiytor/invenda/api/pkg/models"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// User define uma série de operações relacionadas a usuarios.
type User interface {
	// UserGet recupera um usuario com ID `id` do banco de dados. Retorna um ErrNotFound quando nenhum usuario é encotrado.
	UserGet(ctx context.Context, id string) (*models.User, error)

	// UserConflicts verifica se algum campo de `fields` existe no banco de dados. Retorna uma lista dos dados conflitantes
	// ou um erro.
	UserConflicts(ctx context.Context, fields map[string]string) ([]string, error)

	// UserCreate cria um novo usuario no banco de dados. Retorna o ID inserido ou um erro.
	UserCreate(ctx context.Context, user *models.User) (string, error)

	// UserSetConfirmed marca um usuario com ID `id` como confirmado ou não.
	UserSetConfirmed(ctx context.Context, id string, confirmed bool) error
}

func (s *store) UserGet(ctx context.Context, id string) (*models.User, error) {
	user := new(models.User)

	query := s.users.FindOne(ctx, bson.M{"_id": id})
	if err := query.Decode(user); err != nil {
		return nil, fromMongoError(err)
	}

	return user, nil
}

// TODO: fazer fields ser uma struct especializada
func (s *store) UserConflicts(ctx context.Context, fields map[string]string) ([]string, error) {
	pipeline := []bson.M{}

	if email, ok := fields["email"]; ok && email != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"email": fields["email"]}})
	}

	cursor, err := s.users.Aggregate(ctx, pipeline)
	if err != nil {
		return []string{}, fromMongoError(err)
	}
	defer cursor.Close(ctx)

	conflicts := []string{}
	for cursor.Next(ctx) {
		user := new(models.User)
		if err = cursor.Decode(user); err != nil {
			return []string{}, fromMongoError(err)
		}

		if email, ok := fields["email"]; ok && user.Email == email {
			conflicts = append(conflicts, "email")
		}
	}

	return conflicts, nil
}

func (s *store) UserCreate(ctx context.Context, user *models.User) (string, error) {
	// TODO: nós não nos importamos (ainda) com performance
	user.ID = ulid.Make().String()
	_, err := s.users.InsertOne(ctx, user)

	return user.ID, fromMongoError(err)
}

func (s *store) UserSetConfirmed(ctx context.Context, id string, confirmed bool) error {
	res, err := s.users.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"confirmed": confirmed}})
	if err != nil {
		return fromMongoError(err)
	}

	if res.MatchedCount < 1 {
		return ErrNotFound
	}

	return nil
}
