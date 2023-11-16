package codefixture

type Fixture struct {
	models map[ModelRef]any
}

func NewFixture() *Fixture {
	return &Fixture{
		models: make(map[ModelRef]any),
	}
}

func (f *Fixture) GetModel(ref ModelRef) any {
	return f.models[ref]
}

func (f *Fixture) SetModel(ref ModelRef, m any) {
	f.models[ref] = m
}
