package schema

type Hashable interface {
	Identifiable

	Size() int

	Empty() Hashable
}
