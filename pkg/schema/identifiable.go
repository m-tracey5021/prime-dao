package schema

type Identifiable interface {
	Id() uint64

	Descriptor() string
}
