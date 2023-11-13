package codefixture

import (
	"fmt"
	"reflect"
	"sort"
)

type FixtureBuilder struct {
	constructors map[reflect.Type]func() any
	writers      map[reflect.Type]func(m any) (any, error)
	models       map[ModelRef]any
	relations    []ModelRelation
}

func NewFixtureBuilder() *FixtureBuilder {
	return &FixtureBuilder{
		constructors: map[reflect.Type]func() any{},
		writers:      make(map[reflect.Type]func(m any) (any, error)),
		models:       make(map[ModelRef]any),
	}
}

func RegisterWriter[T any](b *FixtureBuilder, f func(m T) (T, error)) error {
	ptrType := reflect.TypeOf((*T)(nil)).Elem()
	if ptrType.Kind() != reflect.Ptr {
		return fmt.Errorf("type %v is not a pointer", ptrType)
	}

	b.writers[ptrType] = func(m any) (any, error) {
		switch m := m.(type) {
		case T:
			return f(m)
		default:
			return nil, fmt.Errorf("invalid type: %T", m)
		}
	}

	return nil
}

func RegisterConstructor[T any](b *FixtureBuilder, constructor func() T) error {
	ptrType := reflect.TypeOf((*T)(nil)).Elem()
	if ptrType.Kind() != reflect.Ptr {
		return fmt.Errorf("type %v is not a pointer", ptrType)
	}

	b.constructors[ptrType] = func() any {
		return constructor()
	}

	return nil
}

func AddModel[T any](b *FixtureBuilder, setter func(T)) (ModelRef, error) {
	ptrType := reflect.TypeOf((*T)(nil)).Elem()
	if ptrType.Kind() != reflect.Ptr {
		return "", fmt.Errorf("type %v is not a pointer", ptrType)
	}

	structType := ptrType.Elem()
	if structType.Kind() != reflect.Struct {
		return "", fmt.Errorf("type %v is not a struct", structType)
	}

	ref := NewModelRef()
	var m any

	contructor := b.constructors[ptrType]
	if contructor == nil {
		m = reflect.New(structType).Interface()
	} else {
		m = contructor()
	}

	if setter != nil {
		switch m := m.(type) {
		case T:
			setter(m)
		default:
			panic(fmt.Sprintf("invalid type: %T", m))
		}
	}

	b.models[ref] = m
	return ref, nil
}

func AddRelation[T any, U any](b *FixtureBuilder, target ModelRef, foreign ModelRef, connector func(T, U)) {
	b.relations = append(b.relations, ModelRelation{
		TargetRef:  target,
		ForeignRef: foreign,
		Connector: func(target, dependent any) {
			connector(target.(T), dependent.(U))
		},
	})
}

func (b *FixtureBuilder) GetModel(ref ModelRef) any {
	return b.models[ref]
}

func (b *FixtureBuilder) Build() (*Fixture, error) {
	f := NewFixture()
	refs, inModels := b.getModelsOrderedByRelations()

	for i, ref := range refs {
		inModel := inModels[i]

		for _, relation := range b.relations {
			if relation.TargetRef != ref {
				continue
			}
			foreignModel := f.GetModel(relation.ForeignRef)
			if foreignModel == nil {
				return nil, fmt.Errorf("model ref %s is not found", relation.ForeignRef)
			}
			relation.Connector(inModel, foreignModel)
		}

		typ := reflect.TypeOf(inModel)
		writer, ok := b.writers[typ]
		if !ok {
			return nil, fmt.Errorf("writer for %T is not found", inModel)
		}

		outModel, err := writer(inModel)
		if err != nil {
			return nil, err
		}
		f.SetModel(ref, outModel)
	}

	return f, nil
}

// getModelsOrderedByRelations returns a slice of ModelRef and a slice of models
// ordered based on their hierarchical depth as defined in Relations.
func (ib *FixtureBuilder) getModelsOrderedByRelations() ([]ModelRef, []any) {
	// Initialize depths for each model in InModels
	depths := make(map[ModelRef]int)
	for ref := range ib.models {
		depths[ref] = 0
	}

	// Update depths based on Relations
	changed := true
	for changed {
		changed = false
		for _, relation := range ib.relations {
			targetDepth := depths[relation.TargetRef]
			foreignDepth := depths[relation.ForeignRef]

			if targetDepth >= foreignDepth {
				newDepth := targetDepth + 1
				if newDepth > depths[relation.ForeignRef] {
					depths[relation.ForeignRef] = newDepth
					changed = true
				}
			}
		}
	}

	// Create a slice for sorting based on depth
	type modelWithDepth struct {
		ref   ModelRef
		depth int
		model any
	}

	var modelsWithDepth []modelWithDepth
	for ref, m := range ib.models {
		modelsWithDepth = append(modelsWithDepth, modelWithDepth{ref, depths[ref], m})
	}

	// Sort models based on depth
	sort.Slice(modelsWithDepth, func(i, j int) bool {
		return modelsWithDepth[i].depth > modelsWithDepth[j].depth
	})

	// Prepare the final sorted results
	var orderedRefs []ModelRef
	var orderedModels []any
	for _, m := range modelsWithDepth {
		orderedRefs = append(orderedRefs, m.ref)
		orderedModels = append(orderedModels, m.model)
	}

	return orderedRefs, orderedModels
}
