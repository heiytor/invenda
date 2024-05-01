package store

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type Entity interface {
	// Entity returns the related entity on which the methods operate.
	Entity() string
}

// or creates a `{ "$match": { "$or": [...] } }` aggregation pipeline. The $match stage
// filters documents to include only those that match any non-zero field of the target object.
func or[T any](target *T) []bson.M {
	conditions := []bson.M{}

	targetValue := reflect.ValueOf(*target)
	if targetValue.Kind() != reflect.Struct {
		return conditions
	}

	targetType := reflect.TypeOf(*target)
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldValue := reflect.ValueOf(target).Elem().FieldByName(field.Name)

		if !reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(fieldValue.Type()).Interface()) {
			conditions = append(conditions, bson.M{field.Tag.Get("bson"): fieldValue.Interface()})
		}
	}

	return []bson.M{
		{
			"$match": bson.M{
				"$or": conditions,
			},
		},
	}
}

// partialEqual reports whether b.N has the same value as a.N where N is all non-zero
// attributes of a.
func partialEqual[T any](a, b *T) []string {
	conflicts := make([]string, 0)

	ra := reflect.ValueOf(*a)
	rb := reflect.ValueOf(*b)

	for i := 0; i < ra.NumField(); i++ {
		name := ra.Type().Field(i).Name

		field1 := ra.Field(i)
		field2 := rb.FieldByName(name)

		if field2.IsValid() && reflect.DeepEqual(field1.Interface(), field2.Interface()) {
			conflicts = append(conflicts, strings.ToLower(name))
		}
	}

	return conflicts
}
