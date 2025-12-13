package valueobject

// Default values and limits for pagination.
const (
	// DefaultPage is the default page number when none is specified.
	DefaultPage = 1
	// DefaultPageSize is the default number of items per page.
	DefaultPageSize = 20
	// MaxPageSize is the maximum allowed items per page to prevent excessive data retrieval.
	MaxPageSize = 100
	// MinPageSize is the minimum allowed items per page.
	MinPageSize = 1
)

// Pagination represents pagination parameters for database queries.
// It is an immutable value object that automatically normalizes values
// to ensure they fall within acceptable ranges.
type Pagination struct {
	page     int
	pageSize int
}

// NewPagination creates a new Pagination instance with validated and normalized values.
// Invalid or out-of-range values are automatically adjusted:
//   - page < 1 is set to DefaultPage (1)
//   - pageSize < MinPageSize is set to DefaultPageSize (20)
//   - pageSize > MaxPageSize is capped at MaxPageSize (100)
func NewPagination(page, pageSize int) Pagination {
	p := Pagination{
		page:     page,
		pageSize: pageSize,
	}

	// Normalize values to valid ranges
	p.normalize()

	return p
}

// DefaultPagination returns a Pagination instance with default values.
// Uses DefaultPage (1) and DefaultPageSize (20).
func DefaultPagination() Pagination {
	return Pagination{
		page:     DefaultPage,
		pageSize: DefaultPageSize,
	}
}

// normalize adjusts pagination values to valid ranges.
// This is an internal method called during construction.
func (p *Pagination) normalize() {
	// Minimum page is 1
	if p.page < 1 {
		p.page = DefaultPage
	}

	// PageSize must be between MinPageSize and MaxPageSize
	if p.pageSize < MinPageSize {
		p.pageSize = DefaultPageSize
	}

	if p.pageSize > MaxPageSize {
		p.pageSize = MaxPageSize
	}
}

// Page returns the current page number (1-indexed).
func (p Pagination) Page() int {
	return p.page
}

// PageSize returns the number of items per page.
func (p Pagination) PageSize() int {
	return p.pageSize
}

// Offset calculates the offset for SQL OFFSET clause.
// Example: Page 3 with PageSize 20 returns Offset 40.
func (p Pagination) Offset() int {
	return (p.page - 1) * p.pageSize
}

// Limit returns the limit for SQL LIMIT clause.
// This is equivalent to PageSize.
func (p Pagination) Limit() int {
	return p.pageSize
}

// HasNextPage determines if there is a next page given the total number of items.
// Returns true if there are more items beyond the current page.
func (p Pagination) HasNextPage(totalItems int64) bool {
	return int64(p.page*p.pageSize) < totalItems
}

// HasPreviousPage determines if there is a previous page.
// Returns true if the current page is greater than 1.
func (p Pagination) HasPreviousPage() bool {
	return p.page > 1
}

// TotalPages calculates the total number of pages given the total number of items.
// Returns 0 if totalItems is 0, otherwise rounds up to ensure all items are accessible.
func (p Pagination) TotalPages(totalItems int64) int {
	if totalItems == 0 {
		return 0
	}

	pages := totalItems / int64(p.pageSize)
	if totalItems%int64(p.pageSize) > 0 {
		pages++
	}

	return int(pages)
}

// PaginatedResult represents the result of a paginated query.
// It is a generic type that can hold any slice of items along with
// pagination metadata for API responses.
type PaginatedResult[T any] struct {
	// Items contains the slice of results for the current page.
	Items []T `json:"items"`
	// TotalItems is the total count of items across all pages.
	TotalItems int64 `json:"total_items"`
	// TotalPages is the total number of pages available.
	TotalPages int `json:"total_pages"`
	// CurrentPage is the current page number (1-indexed).
	CurrentPage int `json:"current_page"`
	// PageSize is the number of items per page.
	PageSize int `json:"page_size"`
	// HasNext indicates whether there is a next page available.
	HasNext bool `json:"has_next"`
	// HasPrevious indicates whether there is a previous page available.
	HasPrevious bool `json:"has_previous"`
}

// NewPaginatedResult creates a new PaginatedResult from items, total count, and pagination parameters.
// It automatically calculates all pagination metadata based on the provided Pagination.
//
// Parameters:
//   - items: the slice of items for the current page
//   - totalItems: the total count of items across all pages
//   - pagination: the Pagination value object with page and pageSize
//
// Returns a fully populated PaginatedResult with all metadata fields.
func NewPaginatedResult[T any](items []T, totalItems int64, pagination Pagination) PaginatedResult[T] {
	return PaginatedResult[T]{
		Items:       items,
		TotalItems:  totalItems,
		TotalPages:  pagination.TotalPages(totalItems),
		CurrentPage: pagination.Page(),
		PageSize:    pagination.PageSize(),
		HasNext:     pagination.HasNextPage(totalItems),
		HasPrevious: pagination.HasPreviousPage(),
	}
}
