package codefixture

import (
	"fmt"
	"reflect"
)

func RegisterWriter[T, U any](b *FixtureBuilder, writer func(m T) (U, error)) error {
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

	ref := NewTypedModelRef[T]()
	err := b.addModel(ptrType, ref.ModelRef(), func(m any) {
		t, ok := m.(T)
		if !ok {
			panic(NewInvalidTypeError(m))
		}

		if setter != nil {
			setter(t)
		}
	})
	if err != nil {
		return "", err
	}

	return ref, nil
}

func AddModelWithRelation[T, U any](b *FixtureBuilder, foreign TypedModelRef[U], connector func(any, any)) (TypedModelRef[T], error) {
	ptrType := reflect.TypeOf((*T)(nil)).Elem()
	if ptrType.Kind() != reflect.Ptr {
		return "", NewNotPointerError(ptrType)
	}

	target, err := AddModel[T](b, nil)
	if err != nil {
		return "", err
	}

	err = AddRelation[T, U](b, target, foreign, connector)
	if err != nil {
		return "", err
	}

	return target, nil
}

func ConvertAndAddModel[T, U any](b *FixtureBuilder, setter func(T)) (TypedModelRef[T], error) {
	ptrType := reflect.TypeOf((*T)(nil)).Elem()
	if ptrType.Kind() != reflect.Ptr {
		return "", NewNotPointerError(ptrType)
	}

	ref := NewTypedModelRef[T]()
	err := b.addModel(ptrType, ref.ModelRef(), func(m any) {
		t, ok := m.(T)
		if !ok {
			panic(NewInvalidTypeError(m))
		}

		if setter != nil {
			setter(t)
		}
	})
	if err != nil {
		return "", err
	}

	return ref, nil
}

func AddRelation[T, U any](b *FixtureBuilder, target TypedModelRef[T], foreign TypedModelRef[U], connector func(any, any)) error {
	return b.addRelation(target.ModelRef(), foreign.ModelRef(), func(target, foreign any) {
		t, ok := target.(T)
		if !ok {
			expectedType := reflect.TypeOf((*T)(nil)).Elem()
			actualType := reflect.TypeOf(target)
			panic(NewUnexpectedTypeError(expectedType, actualType))
		}
		u, ok := foreign.(U)
		if !ok {
			expectedType := reflect.TypeOf((*U)(nil)).Elem()
			actualType := reflect.TypeOf(foreign)
			panic(NewUnexpectedTypeError(expectedType, actualType))
		}

		if connector != nil {
			connector(t, u)
		}
	})
}

func GetBuilderModel[T any](b *FixtureBuilder, ref TypedModelRef[T]) T {
	m := b.GetBuilderModel(ref.ModelRef())

	t, ok := m.(T)
	if !ok {
		expectedType := reflect.TypeOf((*T)(nil)).Elem()
		actualType := reflect.TypeOf(m)
		panic(NewUnexpectedTypeError(expectedType, actualType))
	}

	return t
}
