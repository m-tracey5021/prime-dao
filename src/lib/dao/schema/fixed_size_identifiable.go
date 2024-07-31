package schema

type FixedSizeIdentifiable interface {
	FixedSize

	Identifiable
}
