package query

type Query struct {
	Paginator
	Sorter
}

func New() *Query {
	return &Query{
		Paginator: Paginator{},
		Sorter:    Sorter{},
	}
}
