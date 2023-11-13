package codefixture

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name string
}

func TestRegisterConstructor(t *testing.T) {
	t.Run("set default value", func(t *testing.T) {
		b := NewFixtureBuilder()
		err := RegisterConstructor[*Person](b, func() *Person {
			return &Person{Name: "default"}
		})
		assert.NoError(t, err)

		AddModel[*Person](b, func(p *Person) {
			assert.Equal(t, "default", p.Name)
		})
	})
}

func TestAddModel(t *testing.T) {
	t.Run("no setter, has constructor", func(t *testing.T) {
		b := NewFixtureBuilder()
		err := RegisterConstructor[*Person](b, func() *Person {
			return &Person{Name: "default"}
		})
		assert.NoError(t, err)

		ref, err := AddModel[*Person](b, nil)
		assert.NoError(t, err)
		assert.Equal(t, "default", b.InModels[ref].(*Person).Name)
	})
	t.Run("no setter, no constructor", func(t *testing.T) {
		b := NewFixtureBuilder()
		ref, err := AddModel[*Person](b, nil)
		assert.NoError(t, err)
		assert.Zero(t, b.InModels[ref].(*Person).Name)
	})
	t.Run("override value by setter", func(t *testing.T) {
		b := NewFixtureBuilder()
		err := RegisterConstructor[*Person](b, func() *Person {
			return &Person{Name: "default"}
		})
		assert.NoError(t, err)

		ref, err := AddModel[*Person](b, func(p *Person) {
			p.Name = "override"
		})
		assert.NoError(t, err)

		assert.Equal(t, "override", b.InModels[ref].(*Person).Name)
	})
}

func TestFixtureBuilder_getInModelsOrderedByRelations(t *testing.T) {
	t.Run("no relations", func(t *testing.T) {
		b := NewFixtureBuilder()
		b.InModels["a"] = "a"
		b.InModels["b"] = "b"

		refs, ms := b.getInModelsOrderedByRelations()
		assert.Len(t, refs, 2)
		assert.Len(t, ms, 2)
	})
	t.Run("leaf should appear first", func(t *testing.T) {
		b := NewFixtureBuilder()
		b.InModels["root"] = "root"
		b.InModels["branch"] = "branch"
		b.InModels["leaf"] = "leaf"
		b.Relations = []ModelRelation{
			{
				TargetRef:  "root",
				ForeignRef: "branch",
			},
			{
				TargetRef:  "branch",
				ForeignRef: "leaf",
			},
		}

		refs, ms := b.getInModelsOrderedByRelations()
		assert.Len(t, refs, 3)
		assert.Len(t, ms, 3)
		assert.Equal(t, "leaf", ms[0])
		assert.Equal(t, "branch", ms[1])
		assert.Equal(t, "root", ms[2])
	})
	t.Run("delayed dependency resolution", func(t *testing.T) {
		b := NewFixtureBuilder()
		b.InModels["d"] = "d"
		b.InModels["c"] = "c"
		b.InModels["b"] = "b"
		b.InModels["a"] = "a"

		b.Relations = []ModelRelation{
			{
				TargetRef:  "b",
				ForeignRef: "a",
			},
			{
				TargetRef:  "c",
				ForeignRef: "b",
			},
			{
				TargetRef:  "d",
				ForeignRef: "a",
			},
			{
				TargetRef:  "d",
				ForeignRef: "c",
			},
		}

		refs, ms := b.getInModelsOrderedByRelations()
		assert.Len(t, refs, 4)
		assert.Len(t, ms, 4)
		assert.Equal(t, []any{"a", "b", "c", "d"}, ms)
	})
}
