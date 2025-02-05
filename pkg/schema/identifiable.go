package schema

import "github.com/google/uuid"

type Identifiable interface {
	Id() uuid.UUID

	SetId(uuid.UUID)

	Descriptor() string
}
