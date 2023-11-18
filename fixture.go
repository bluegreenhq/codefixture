package codefixture

import "log"

type Fixture struct {
	models map[ModelRef]any
}

func NewFixture() *Fixture {
	return &Fixture{
		models: make(map[ModelRef]any),
	}
}

func (f *Fixture) GetModel(ref ModelRef) any {
	m, ok := f.models[ref]
	if !ok {
		panic(NewModelRefNotFoundError(ref))
	}
	return m
}

func (f *Fixture) SetModel(ref ModelRef, m any) {
	log.Printf("Fixture.SetModel model=%T", m)
	f.models[ref] = m
}
