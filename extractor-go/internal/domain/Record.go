package domain

type Record[T any] struct {
	Update string
	Data   T
}
