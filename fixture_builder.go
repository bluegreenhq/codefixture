package codefixture

import (
	"fmt"
	"reflect"
	"sort"
)

type FixtureBuilder struct {
	Constructors map[reflect.Type]func() any
	Writers      map[reflect.Type]func(m any) (any, error)
	InModels     map[ModelRef]any
	OutModels    map[ModelRef]any
	Relations    []ModelRelation
}

func NewFixtureBuilder() *FixtureBuilder {
	return &FixtureBuilder{
		Constructors: map[reflect.Type]func() any{},
		Writers:      make(map[reflect.Type]func(m any) (any, error)),
		InModels:     make(map[ModelRef]any),
		OutModels:    make(map[ModelRef]any),
	}
}

func RegisterWriter[T any](b *FixtureBuilder, f func(m T) (T, error)) error {
	ptrType := reflect.TypeOf((*T)(nil)).Elem()
	if ptrType.Kind() != reflect.Ptr {
		return fmt.Errorf("type %v is not a pointer", ptrType)
	}

	b.Writers[ptrType] = func(m any) (any, error) {
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

	b.Constructors[ptrType] = func() any {
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

	contructor := b.Constructors[ptrType]
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

	b.InModels[ref] = m
	return ref, nil
}

func AddRelation[T any, U any](b *FixtureBuilder, target ModelRef, foreign ModelRef, connector func(T, U)) {
	b.Relations = append(b.Relations, ModelRelation{
		TargetRef:  target,
		ForeignRef: foreign,
		Connector: func(target, dependent any) {
			connector(target.(T), dependent.(U))
		},
	})
}

func (b *FixtureBuilder) Build() error {
	refs, inModels := b.getInModelsOrderedByRelations()

	for i, modelAlias := range refs {
		inModel := inModels[i]

		for _, relation := range b.Relations {
			if relation.TargetRef != modelAlias {
				continue
			}
			dstModel, ok := b.OutModels[relation.ForeignRef]
			if !ok {
				return fmt.Errorf("model alias %s is not found", relation.ForeignRef)
			}
			relation.Connector(inModel, dstModel)
		}

		typ := reflect.TypeOf(inModel)
		writer, ok := b.Writers[typ]
		if !ok {
			return fmt.Errorf("writer for %T is not found", inModel)
		}

		outModel, err := writer(inModel)
		if err != nil {
			return err
		}
		b.OutModels[modelAlias] = outModel
	}

	return nil
}

func GetModels[T any](b *FixtureBuilder) []T {
	var models []T
	for _, m := range b.OutModels {
		switch m := m.(type) {
		case T:
			models = append(models, m)
		}
	}
	return models
}

// getInModelsOrderedByRelations returns a slice of ModelRef and a slice of models
// ordered based on their hierarchical depth as defined in Relations.
func (ib *FixtureBuilder) getInModelsOrderedByRelations() ([]ModelRef, []any) {
	// Initialize depths for each model in InModels
	depths := make(map[ModelRef]int)
	for ref := range ib.InModels {
		depths[ref] = 0
	}

	// Update depths based on Relations
	for _, relation := range ib.Relations {
		targetDepth, targetExists := depths[relation.TargetRef]
		foreignDepth, foreignExists := depths[relation.ForeignRef]

		if targetExists && foreignExists {
			if targetDepth >= foreignDepth {
				depths[relation.ForeignRef] = targetDepth + 1
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
	for ref, model := range ib.InModels {
		modelsWithDepth = append(modelsWithDepth, modelWithDepth{ref, depths[ref], model})
	}

	// Sort models based on depth
	sort.Slice(modelsWithDepth, func(i, j int) bool {
		return modelsWithDepth[i].depth > modelsWithDepth[j].depth
	})

	// Prepare the final sorted results
	var orderedRefs []ModelRef
	var orderedModels []any
	for _, mwd := range modelsWithDepth {
		orderedRefs = append(orderedRefs, mwd.ref)
		orderedModels = append(orderedModels, mwd.model)
	}

	return orderedRefs, orderedModels
}
