package codefixture

import (
	"log"
	"reflect"
	"sort"
)

type FixtureBuilder struct {
	constructors map[reflect.Type]Constructor
	writers      map[reflect.Type]Writer
	models       map[ModelRef]any
	converted    map[ModelRef]ModelRef
	relations    []ModelRelation
	option       *FixtureBuilderOption
}
type FixtureBuilderOption struct {
	AllowEmptyWriter bool
}

type Constructor func() any
type Setter func(any)
type Writer func(m any) (any, error)

func NewFixtureBuilder() *FixtureBuilder {
	return &FixtureBuilder{
		constructors: make(map[reflect.Type]Constructor),
		writers:      make(map[reflect.Type]Writer),
		models:       make(map[ModelRef]any),
		converted:    make(map[ModelRef]ModelRef),
		option:       &FixtureBuilderOption{},
	}
}

func NewFixtureBuilderWithOption(option *FixtureBuilderOption) *FixtureBuilder {
	return &FixtureBuilder{
		constructors: make(map[reflect.Type]Constructor),
		writers:      make(map[reflect.Type]Writer),
		models:       make(map[ModelRef]any),
		converted:    make(map[ModelRef]ModelRef),
		option:       option,
	}
}

func (b *FixtureBuilder) RegisterWriter(typeInstance any, writer Writer) error {
	ptrType := reflect.TypeOf(typeInstance)

	return b.registerWriter(ptrType, writer)
}

func (b *FixtureBuilder) RegisterConstructor(typeInstance any, constructor Constructor) error {
	ptrType := reflect.TypeOf(typeInstance)

	return b.registerConstructor(ptrType, constructor)
}

func (b *FixtureBuilder) AddModel(m any) (ModelRef, error) {
	ptrType := reflect.TypeOf(m)
	if ptrType.Kind() != reflect.Ptr {
		return "", NewNotPointerError(ptrType)
	}

	ref := NewModelRef()
	b.SetBuilderModel(ref, m)
	return ref, nil
}

func (b *FixtureBuilder) WithModel(m any, ref ModelRef) *FixtureBuilder {
	ptrType := reflect.TypeOf(m)
	if ptrType.Kind() != reflect.Ptr {
		err := NewNotPointerError(ptrType)
		panic(err)
	}

	b.SetBuilderModel(ref, m)
	return b
}

func (b *FixtureBuilder) AddModelBySetter(typeInstance any, setter Setter) (ModelRef, error) {
	ptrType := reflect.TypeOf(typeInstance)
	if ptrType.Kind() != reflect.Ptr {
		return "", NewNotPointerError(ptrType)
	}

	ref := NewModelRef()
	err := b.addModel(ptrType, ref, setter)
	if err != nil {
		return "", err
	}

	return ref, nil
}

func (b *FixtureBuilder) AddRelation(target ModelRef, foreign ModelRef, connector Connector) error {
	return b.addRelation(target, foreign, connector)
}

func (b *FixtureBuilder) WithRelation(target ModelRef, foreign ModelRef, connector Connector) *FixtureBuilder {
	err := b.addRelation(target, foreign, connector)
	if err != nil {
		panic(err)
	}
	return b
}

func (b *FixtureBuilder) GetBuilderModel(ref ModelRef) any {
	return b.models[ref]
}

func (b *FixtureBuilder) SetBuilderModel(ref ModelRef, m any) {
	b.models[ref] = m
}

func (b *FixtureBuilder) GetBuilderModelResolvingConverted(ref ModelRef) any {
	return b.GetBuilderModel(b.GetModelRefResolvingConverted(ref))
}

func (b *FixtureBuilder) GetModelRefResolvingConverted(ref ModelRef) ModelRef {
	if convertedRef, ok := b.converted[ref]; ok {
		return convertedRef
	}
	return ref
}

func (b *FixtureBuilder) Build() (*Fixture, error) {
	log.Println("FixtureBuilder.Build begin")

	f := NewFixture()
	refs, inModels := b.getModelsOrderedByRelations()

	for i, ref := range refs {
		inModel := inModels[i]

		inType := reflect.TypeOf(inModel)
		writer := b.writers[inType]

		model, err := b.resolveRelations(ref, inModel, f)
		if err != nil {
			return nil, err
		}

		if writer != nil {
			outModel, err := writer(model)
			if err != nil {
				return nil, err
			}
			model = outModel
		} else {
			if !b.option.AllowEmptyWriter {
				return nil, NewWriterNotFoundError(inType)
			}
		}

		f.SetModel(ref, model)

		outType := reflect.TypeOf(model)
		if outType != inType {
			outRef := NewModelRef()
			b.converted[ref] = outRef
		}
	}

	log.Println("FixtureBuilder.Build end")
	return f, nil
}

func (b *FixtureBuilder) registerWriter(ptrType reflect.Type, writer Writer) error {
	if ptrType.Kind() != reflect.Ptr {
		return NewNotPointerError(ptrType)
	}

	b.writers[ptrType] = writer
	return nil
}

func (b *FixtureBuilder) registerConstructor(ptrType reflect.Type, constructor Constructor) error {
	if ptrType.Kind() != reflect.Ptr {
		return NewNotPointerError(ptrType)
	}

	b.constructors[ptrType] = constructor
	return nil
}

func (b *FixtureBuilder) addModel(ptrType reflect.Type, ref ModelRef, setter Setter) error {
	structType := ptrType.Elem()
	if structType.Kind() != reflect.Struct {
		return NewNotStructError(structType)
	}

	var m any
	contructor := b.constructors[ptrType]
	if contructor == nil {
		m = reflect.New(structType).Interface()
	} else {
		m = contructor()
	}

	if setter != nil {
		setter(m)
	}

	b.SetBuilderModel(ref, m)
	return nil
}

func (b *FixtureBuilder) addRelation(target ModelRef, foreign ModelRef, connector Connector) error {
	targetModel := b.GetBuilderModelResolvingConverted(target)
	if targetModel == nil {
		return NewModelRefNotFoundError(target)
	}
	foreignModel := b.GetBuilderModelResolvingConverted(foreign)
	if foreignModel == nil {
		return NewModelRefNotFoundError(foreign)
	}

	b.relations = append(b.relations, ModelRelation{
		TargetRef:  target,
		ForeignRef: foreign,
		Connector:  connector,
	})
	return nil
}

func (b *FixtureBuilder) resolveRelations(ref ModelRef, model any, f *Fixture) (any, error) {
	log.Printf("FixtureBuilder.resolveRelations model=%T", model)
	var targetModel = model

	for _, relation := range b.relations {
		targetRef := b.GetModelRefResolvingConverted(relation.TargetRef)

		if targetRef != ref {
			continue
		}

		targetModel = b.GetBuilderModelResolvingConverted(ref)
		foreignModel := f.GetModel(b.GetModelRefResolvingConverted(relation.ForeignRef))
		if foreignModel == nil {
			return nil, NewModelRefNotFoundError(relation.ForeignRef)
		}
		log.Printf("FixtureBuilder.resolveRelations target=%T, foreign=%T\n", targetModel, foreignModel)
		relation.Connector(targetModel, foreignModel)
	}

	return targetModel, nil
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
