package internal

import (
	"github.com/heiytor/invenda/api/pkg/query"
	"go.mongodb.org/mongo-driver/bson"
)

// FromPaginator converts the Paginator instance to a BSON pagination expression for MongoDB queries.
// If the per-page count is less than 1, it returns nil.
func FromPaginator(p *query.Paginator) []bson.M {
	if p == nil {
		return []bson.M{}
	}

	return []bson.M{
		{"$skip": p.Size * (p.Page - 1)},
		{"$limit": p.Size},
	}
}

// FromSorter converts the Sort instance to a BSON sorting expression for MongoDB queries.
// If an invalid value of `Sorter.By` is provided, it defaults to ascending order (OrderAsc).
func FromSorter(s *query.Sorter) []bson.M {
	options := map[string]int{
		query.OrderAsc:  1,
		query.OrderDesc: -1,
	}

	order, ok := options[s.Order]
	if !ok {
		order = -1
	}

	return []bson.M{
		{
			"$sort": bson.M{
				s.By: order,
			},
		},
	}
}
