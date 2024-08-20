package schema

type Identifiable interface {
	Id() uint64

	SetId(uint64) Identifiable
}
