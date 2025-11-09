package domain

type Identifiable[K comparable] interface {
	GetId() K
}

type IdentifiablePointer[K comparable, T any] interface {
	SetId(id K)
	*T
}
