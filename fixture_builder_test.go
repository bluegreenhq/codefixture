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

func TestFixtureBuilder_AddRelation(t *testing.T) {
	t.Run("add relation", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()

		p, _ := b.AddModel(&Person{}, func(m any) {
			m.(*Person).ID = 1
		})
		g, _ := b.AddModel(&Group{}, func(m any) {
			m.(*Group).ID = 2
		})

		b.AddRelation(p, g, func(p, g any) {
			p.(*Person).GroupID = g.(*Group).ID
		})

		f, err := b.Build()
		assert.NoError(t, err)
		assert.Equal(t, 2, f.GetModel(p).(*Person).GroupID)
	})
}
