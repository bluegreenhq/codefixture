package codefixture

func GetModel[T any](f *Fixture, ref ModelRef) T {
	m := f.GetModel(ref)

	t, ok := m.(T)
	if !ok {
		panic(NewInvalidTypeError(m))
	}

	return t
}
