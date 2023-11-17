package codefixture_test

import (
	"testing"

	"github.com/bluegreenhq/codefixture"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	ID      int
	GroupID int
	Name    string
}

type Group struct {
	ID   int
	Name string
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

		p, err := codefixture.ConvertAndAddModel[*PersonMaterial, *Person](b, nil)
		assert.NoError(t, err)
		g, err := codefixture.AddModel(b, func(g *Group) {
			g.Name = "family"
		})
		assert.NoError(t, err)

		err = codefixture.AddRelation(b, p, g, func(p *Person, g *Group) {
			p.GroupID = g.ID
		})
		assert.NoError(t, err)

		_, err = b.Build()
		assert.NoError(t, err)
	})
}
