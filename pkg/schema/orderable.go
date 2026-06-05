package schema

type Order int

const (
	Smaller Order = iota

	Equal

	Larger
)

type Orderable interface {
	FixedSizeIdentifiable

	Compare(Orderable) Order
}
