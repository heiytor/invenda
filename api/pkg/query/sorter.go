package query

const (
	OrderAsc  = "asc"
	OrderDesc = "desc"
)

// Sorter represents the sorting order in a query.
type Sorter struct {
	By    string `query:"sort"`
	Order string `query:"order" validate:"in:asc,desc"`
}

// NormalizeWith ensures that the sorting order is valid. If an invalid order is provided, it defaults
// to descending order. The [Sorter.By] is set to by when empty.
func (s *Sorter) NormalizeWith(by string) {
	if s.By == "" {
		s.By = by
	}

	if s.Order != OrderAsc && s.Order != OrderDesc {
		s.Order = OrderDesc
	}
}
