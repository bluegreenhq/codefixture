package codefixture

import (
	"fmt"
	"reflect"
)

func RegisterWriter[T any, U any](b *FixtureBuilder, writer func(m T) (U, error)) error {
	ptrType := reflect.TypeOf((*T)(nil)).Elem()

	return b.registerWriter(ptrType, func(m any) (any, error) {
		t, ok := m.(T)
		if !ok {
			return nil, fmt.Errorf("invalid type: %T", m)
		}
		return writer(t)
	})
}

func RegisterConstructor[T any](b *FixtureBuilder, constructor func() T) error {
	ptrType := reflect.TypeOf((*T)(nil)).Elem()

	b.registerConstructor(ptrType, func() any {
		return constructor()
	})
	return nil
}

func AddModel[T any](b *FixtureBuilder, setter func(T)) (TypedModelRef[T], error) {
	ptrType := reflect.TypeOf((*T)(nil)).Elem()
	if ptrType.Kind() != reflect.Ptr {
		return "", NewNotPointerError(ptrType)
	}

	ref, err := b.addModel(ptrType, func(m any) {
		t, ok := m.(T)
		if !ok {
			panic(NewInvalidTypeError(m))
		}

		if setter != nil {
			setter(t)
		}
	})

	return TypedModelRef[T](ref), err
}

func AddRelation[T any, U any](b *FixtureBuilder, target TypedModelRef[T], foreign TypedModelRef[U], connector func(T, U)) error {
	return b.addRelation(target.ModelRef(), foreign.ModelRef(), func(target, foreign any) {
		t, ok := target.(T)
		if !ok {
			panic(NewInvalidTypeError(target))
		}
		u, ok := foreign.(U)
		if !ok {
			panic(NewInvalidTypeError(foreign))
		}

		if connector != nil {
			connector(t, u)
		}
	})
}

func GetBuilderModel[T any](b *FixtureBuilder, ref TypedModelRef[T]) T {
	m := b.models[ref.ModelRef()]

	t, ok := m.(T)
	if !ok {
		panic(NewInvalidTypeError(m))
	}

	return t
}
