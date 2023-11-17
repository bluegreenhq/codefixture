package codefixture

func GetModel[T any](f *Fixture, ref TypedModelRef[T]) T {
	m, ok := f.models[ref.ModelRef()]
	if !ok {
		panic(NewModelRefNotFoundError(ref.ModelRef()))
	}

	t, ok := m.(T)
	if !ok {
		panic(NewInvalidTypeError(m))
	}

	return t
}

func GetModels[T any](f *Fixture) []T {
	var models []T
	for _, m := range f.models {
		switch m := m.(type) {
		case T:
			models = append(models, m)
		}
	}
	return models
}
