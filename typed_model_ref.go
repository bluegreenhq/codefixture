package codefixture

type TypedModelRef[T any] ModelRef

func NewTypedModelRef[T any]() TypedModelRef[T] {
	return TypedModelRef[T](NewModelRef())
}

func (r TypedModelRef[T]) ModelRef() ModelRef {
	return ModelRef(r)
}
