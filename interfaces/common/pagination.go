package common

import (
	"math"
)

type OrderBy struct {
	Column    string // Column name
	Direction string // ASC or DESC
}

func NewOrderBy(column, direction string) *OrderBy {
	return &OrderBy{
		Column:    column,
		Direction: direction,
	}
}

type Pagination struct {
	Page      int
	Limit     int
	Offset    int
	OrderBy   []OrderBy
	TotalPage int
}

func NewPagination(page, limit int, orderBy []OrderBy) *Pagination {
	offset := (page - 1) * limit

	return &Pagination{
		Page:    page,
		Limit:   limit,
		Offset:  offset,
		OrderBy: orderBy,
	}
}

func (p *Pagination) SetTotalPageByTotalItems(total int) {
	p.TotalPage = int(math.Ceil(float64(total) / float64(p.Limit)))
}

func (p *Pagination) GetNextPage() int {
	if p.Page < p.TotalPage {
		return p.Page + 1
	}
	return 0
}
