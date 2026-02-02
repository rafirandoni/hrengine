package common

import "math"

type Pagination struct {
	TotalItem   int
	TotalPage   int
	CurrentPage int
	ItemPerPage int
}

func MakePagination(currentPage int, itemPerPage int, totalItem int) *Pagination {
	totalPage := math.Ceil(float64(totalItem) / float64(itemPerPage))

	return &Pagination{
		TotalItem:   totalItem,
		TotalPage:   int(totalPage),
		CurrentPage: currentPage,
		ItemPerPage: itemPerPage,
	}
}
