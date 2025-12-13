package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/daniel-caso-github/realtime-alerting-system/internal/domain/valueobject"
)

func TestNewPagination_Normalization(t *testing.T) {
	testCases := []struct {
		name             string
		inputPage        int
		inputPageSize    int
		expectedPage     int
		expectedPageSize int
	}{
		{"valid values", 2, 20, 2, 20},
		{"negative page", -1, 20, 1, 20},
		{"zero page", 0, 20, 1, 20},
		{"negative page size", 2, -5, 2, valueobject.DefaultPageSize},
		{"zero page size", 2, 0, 2, valueobject.DefaultPageSize},
		{"page size too large", 2, 500, 2, valueobject.MaxPageSize},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := valueobject.NewPagination(tc.inputPage, tc.inputPageSize)

			assert.Equal(t, tc.expectedPage, p.Page())
			assert.Equal(t, tc.expectedPageSize, p.PageSize())
		})
	}
}

func TestPagination_Offset(t *testing.T) {
	testCases := []struct {
		page           int
		pageSize       int
		expectedOffset int
	}{
		{1, 20, 0},
		{2, 20, 20},
		{3, 20, 40},
		{1, 50, 0},
		{3, 50, 100},
	}

	for _, tc := range testCases {
		p := valueobject.NewPagination(tc.page, tc.pageSize)
		assert.Equal(t, tc.expectedOffset, p.Offset())
	}
}

func TestPagination_TotalPages(t *testing.T) {
	testCases := []struct {
		pageSize   int
		totalItems int64
		expected   int
	}{
		{20, 0, 0},
		{20, 10, 1},
		{20, 20, 1},
		{20, 21, 2},
		{20, 100, 5},
		{20, 101, 6},
	}

	for _, tc := range testCases {
		p := valueobject.NewPagination(1, tc.pageSize)
		assert.Equal(t, tc.expected, p.TotalPages(tc.totalItems))
	}
}

func TestPagination_HasNextPage(t *testing.T) {
	p := valueobject.NewPagination(1, 20)

	assert.False(t, p.HasNextPage(10))
	assert.False(t, p.HasNextPage(20))
	assert.True(t, p.HasNextPage(21))

	p2 := valueobject.NewPagination(2, 20)
	assert.False(t, p2.HasNextPage(40))
	assert.True(t, p2.HasNextPage(50))
}

func TestPagination_HasPreviousPage(t *testing.T) {
	p1 := valueobject.NewPagination(1, 20)
	p2 := valueobject.NewPagination(2, 20)

	assert.False(t, p1.HasPreviousPage())
	assert.True(t, p2.HasPreviousPage())
}

func TestDefaultPagination(t *testing.T) {
	p := valueobject.DefaultPagination()

	assert.Equal(t, valueobject.DefaultPage, p.Page())
	assert.Equal(t, valueobject.DefaultPageSize, p.PageSize())
}

func TestNewPaginatedResult(t *testing.T) {
	items := []string{"a", "b", "c"}
	pagination := valueobject.NewPagination(1, 20)

	result := valueobject.NewPaginatedResult(items, 50, pagination)

	assert.Equal(t, items, result.Items)
	assert.Equal(t, int64(50), result.TotalItems)
	assert.Equal(t, 3, result.TotalPages)
	assert.Equal(t, 1, result.CurrentPage)
	assert.Equal(t, 20, result.PageSize)
	assert.True(t, result.HasNext)
	assert.False(t, result.HasPrevious)
}
