package codefixture

import "log"

type Fixture interface {
	GetModel(ref ModelRef) any
	SetModel(ref ModelRef, m any)
}

type fixture struct {
	models map[ModelRef]any
}

var _ Fixture = (*fixture)(nil)

func NewFixture() Fixture {
	return &fixture{
		models: make(map[ModelRef]any),
	}
}

func (f *fixture) GetModel(ref ModelRef) any {
	m, ok := f.models[ref]
	if !ok {
		panic(NewModelRefNotFoundError(ref))
	}
	return m
}

func (f *fixture) SetModel(ref ModelRef, m any) {
	log.Printf("Fixture.SetModel model=%T", m)
	f.models[ref] = m
}
