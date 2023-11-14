package codefixture

import (
	"fmt"
	"reflect"
)

func RegisterWriter[T any](b *FixtureBuilder, writer func(m T) (T, error)) error {
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

func AddModel[T any](b *FixtureBuilder, setter func(T)) (ModelRef, error) {
	ptrType := reflect.TypeOf((*T)(nil)).Elem()
	if ptrType.Kind() != reflect.Ptr {
		return "", fmt.Errorf("type %v is not a pointer", ptrType)
	}

	return b.addModel(ptrType, func(m any) {
		t, ok := m.(T)
		if !ok {
			panic(fmt.Sprintf("invalid type: %T", m))
		}

		if setter != nil {
			setter(t)
		}
	})
}

func AddRelation[T any, U any](b *FixtureBuilder, target ModelRef, foreign ModelRef, connector func(T, U)) error {
	return b.addRelation(target, foreign, func(target, foreign any) {
		t, ok := target.(T)
		if !ok {
			panic(fmt.Sprintf("invalid type: %T", target))
		}
		u, ok := foreign.(U)
		if !ok {
			panic(fmt.Sprintf("invalid type: %T", target))
		}

		if connector != nil {
			connector(t, u)
		}
	})
}
