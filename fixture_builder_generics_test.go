package codefixture_test

import (
	"testing"

	"github.com/bluegreenhq/codefixture"
	"github.com/stretchr/testify/assert"
)

func TestRegisterConstructor(t *testing.T) {
	t.Run("set default value", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		err := codefixture.RegisterConstructor(b, func() *Person {
			return &Person{Name: "default"}
		})
		assert.NoError(t, err)

		codefixture.AddModel(b, func(p *Person) {
			assert.Equal(t, "default", p.Name)
		})
	})
}

func TestAddModel(t *testing.T) {
	t.Run("no setter, has constructor", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		err := codefixture.RegisterConstructor(b, func() *Person {
			return &Person{Name: "default"}
		})
		assert.NoError(t, err)

		ref, err := codefixture.AddModel[*Person](b, nil)
		assert.NoError(t, err)

		m := codefixture.GetModel[*Person](b, ref)
		assert.Equal(t, "default", m.Name)
	})
	t.Run("no setter, no constructor", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		ref, err := codefixture.AddModel[*Person](b, nil)
		assert.NoError(t, err)

		m := codefixture.GetModel[*Person](b, ref)
		assert.Zero(t, m.Name)
	})
	t.Run("override value by setter", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		err := codefixture.RegisterConstructor(b, func() *Person {
			return &Person{Name: "default"}
		})
		assert.NoError(t, err)

		ref, err := codefixture.AddModel(b, func(p *Person) {
			p.Name = "override"
		})
		assert.NoError(t, err)

		m := codefixture.GetModel[*Person](b, ref)
		assert.Equal(t, "override", m.Name)
	})
}
