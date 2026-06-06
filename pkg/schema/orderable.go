package schema

import "golang.org/x/exp/constraints"

type Order int

const (
	Smaller Order = iota

	Equal

	Larger
)

type Orderable[T constraints.Ordered] interface {
	Identifiable

	SortKeyValue() T

	Compare(Orderable[T]) Order
}
