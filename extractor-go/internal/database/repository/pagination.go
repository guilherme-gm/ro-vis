package repository

import "math"

type Pagination struct {
	Offset int32
	Limit  int32
}

var PaginateAll = Pagination{
	Offset: 0,
	Limit:  math.MaxInt32,
}
