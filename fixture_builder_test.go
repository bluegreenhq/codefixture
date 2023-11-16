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

func TestFixtureBuilder_WithMode(t *testing.T) {
	t.Run("add multiple models", func(t *testing.T) {
		f, err := codefixture.NewFixtureBuilder().
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
