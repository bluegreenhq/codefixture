package codefixture

import "github.com/google/uuid"

type ModelRef string

func NewModelRef() ModelRef {
	return ModelRef(uuid.New().String())
}
