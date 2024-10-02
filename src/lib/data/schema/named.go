package schema

type Named interface {
	DescribedIdentifiable

	Name() string
}
