package query

import "math"

const (
	MinPage        = 1   // MinPage represents the minimum allowed value for the pagination query's Page parameter.
	MinPerPage     = 1   // MinPerPage represents the minimum allowed value for the pagination query's Size parameter.
	DefaultPerPage = 10  // DefaultPerPage represents the default value for the pagination query's Size parameter.
	MaxPerPage     = 100 // MaxPerPage represents the maximum allowed value for the pagination query's Size parameter.
)

// Paginator represents the paginator parameters in a query.
type Paginator struct {
	Page uint `query:"page" validate:"min:1"`         // Page represents the current page number.
	Size uint `query:"size" validate:"min:1|max:100"` // Size represents the number of items per page.
}

// Normalize ensures valid values for Page and Size in the pagination query.
// If query.Size is less than one, it is set to `DefaultPerPage`.
// If query.Page is less than one, it is set to `MinPage`.
// The maximum allowed value for query.Size is `MaxPerPage`.
func (p *Paginator) Normalize() {
	p.Size = uint(math.Max(math.Min(float64(p.Size), float64(MaxPerPage)), float64(DefaultPerPage)))
	p.Page = uint(math.Max(float64(MinPage), float64(p.Page)))
}
