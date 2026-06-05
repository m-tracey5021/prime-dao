package schema

type Order int

const (
	Smaller Order = iota

	Equal

	Larger
)

type Orderable interface {
	Identifiable

	SortKeyValue() any

	Compare(Orderable) Order
}
