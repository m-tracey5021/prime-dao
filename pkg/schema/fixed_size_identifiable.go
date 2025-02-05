package schema

import "github.com/google/uuid"

type FixedSizeIdentifiable interface {
	FixedSize

	Id() uuid.UUID

	Descriptor() string
}
