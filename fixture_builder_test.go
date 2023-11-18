package codefixture_test

import (
	"fmt"
	"testing"

	"github.com/bluegreenhq/codefixture"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	ID      int
	GroupID int
	Name    string
}

type PersonMaterial struct {
	Name    string
	GroupID int
}

type Group struct {
	ID        int
	CreatorID int
	Name      string
}

func TestFixtureBuilder_ReadmeScenario(t *testing.T) {
	// 1. Import codefixture in your test file.
	// 2. Create a `FixtureBuilder`.
	builder := codefixture.NewFixtureBuilder()

	// 3. Register writers with `FixtureBuilder`.
	err := builder.RegisterWriter(&Person{}, func(m any) (any, error) {
		return m, nil
	})
	assert.NoError(t, err)
	err = builder.RegisterWriter(&Group{}, func(m any) (any, error) {
		return m, nil
	})
	assert.NoError(t, err)

	// 4. Add models and relations to `FixtureBuilder`.
	p, _ := builder.AddModel(&Person{Name: "John"})
	g, _ := builder.AddModel(&Group{Name: "Family"})

	builder.AddRelation(p, g, func(p, g any) {
		p.(*Person).GroupID = g.(*Group).ID
	})

	// 5. Build `Fixture` from `FixtureBuilder`.

	fixture, err := builder.Build()
	assert.NoError(t, err)

	// 6. Access your models from `Fixture`.
	fmt.Printf("Person name: %s\n", fixture.GetModel(p).(*Person).Name)
	fmt.Printf("Group name: %s\n", fixture.GetModel(g).(*Group).Name)
}

func TestFixtureBuilder_RegisterWriter(t *testing.T) {
	t.Run("do nothing", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		err := b.RegisterWriter(&Person{}, func(p any) (any, error) {
			return p, nil
		})
		assert.NoError(t, err)
	})
}

func TestFixtureBuilder_WithModel(t *testing.T) {
	t.Run("add multiple models", func(t *testing.T) {
		option := &codefixture.FixtureBuilderOption{
			AllowEmptyWriter: true,
		}
		f, err := codefixture.NewFixtureBuilderWithOption(option).
			WithModel(&Person{Name: "john"}, "p1").
			WithModel(&Group{Name: "family"}, "g1").
			Build()

		assert.NoError(t, err)
		assert.Equal(t, "john", f.GetModel("p1").(*Person).Name)
		assert.Equal(t, "family", f.GetModel("g1").(*Group).Name)
	})
}

func TestFixtureBuilder_AddRelation(t *testing.T) {
	t.Run("add relation", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		err := b.RegisterWriter(&Person{}, func(p any) (any, error) {
			return p, nil
		})
		assert.NoError(t, err)
		err = b.RegisterWriter(&Group{}, func(g any) (any, error) {
			return g, nil
		})
		assert.NoError(t, err)

		p, _ := b.AddModel(&Person{ID: 1})
		g, _ := b.AddModel(&Group{ID: 2})

		b.AddRelation(p, g, func(p, g any) {
			p.(*Person).GroupID = g.(*Group).ID
		})

		f, err := b.Build()
		assert.NoError(t, err)
		assert.Equal(t, 2, f.GetModel(p).(*Person).GroupID)
	})
}

func TestFixtureBuilder_AddModelAndRelation(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		b.RegisterWriter(&Person{}, func(p any) (any, error) {
			return p, nil
		})
		b.RegisterWriter(&Group{}, func(g any) (any, error) {
			return g, nil
		})
		err := b.RegisterConstructor(&Person{}, func() any {
			return &Person{Name: "default"}
		})
		assert.NoError(t, err)

		g, err := b.AddModelBySetter(&Group{}, func(g any) {
			g.(*Group).ID = 1
			g.(*Group).Name = "family"
		})
		assert.NoError(t, err)

		p, err := b.AddModelAndRelation(&Person{}, g, func(p any, g any) {
			p.(*Person).GroupID = g.(*Group).ID
		})
		assert.NoError(t, err)

		f, err := b.Build()
		assert.NoError(t, err)

		m := f.GetModel(p)
		assert.Equal(t, 1, m.(*Person).GroupID)
	})
	t.Run("convert model", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		b.RegisterWriter(&PersonMaterial{}, func(p any) (any, error) {
			return &Person{}, nil
		})
		b.RegisterWriter(&Group{}, func(g any) (any, error) {
			return g, nil
		})

		p, err := b.AddModel(&PersonMaterial{Name: "john"})
		assert.NoError(t, err)
		g, err := b.AddModelAndRelation(&Group{}, p, func(g any, p any) {
			g.(*Group).CreatorID = p.(*Person).ID
		})
		assert.NoError(t, err)
		group := b.GetBuilderModel(g)
		assert.Equal(t, 0, group.(*Group).CreatorID)
	})
}

func TestFixtureBuilder_Build(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		err := b.RegisterWriter(&PersonMaterial{}, func(p any) (any, error) {
			return &Person{
				Name: p.(*PersonMaterial).Name,
			}, nil
		})
		assert.NoError(t, err)
		err = b.RegisterWriter(&Group{}, func(g any) (any, error) {
			return g, nil
		})
		assert.NoError(t, err)

		p, err := b.AddModel(&PersonMaterial{
			Name: "john",
		})
		assert.NoError(t, err)

		{
			person := b.GetBuilderModel(p).(*PersonMaterial)
			assert.Equal(t, "john", person.Name)
		}

		g, err := b.AddModel(&Group{
			Name: "family",
		})

		assert.NoError(t, err)

		err = b.AddRelation(p, g, func(p any, g any) {
			p.(*PersonMaterial).GroupID = g.(*Group).ID
		})
		assert.NoError(t, err)

		fixture, err := b.Build()
		assert.NoError(t, err)

		{
			person := fixture.GetModel(p).(*Person)
			assert.Equal(t, "john", person.Name)
		}
	})
}
