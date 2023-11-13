package codefixture_test

import (
	"testing"

	"github.com/bluegreenhq/codefixture"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name string
}

func TestRegisterConstructor(t *testing.T) {
	t.Run("set default value", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		err := codefixture.RegisterConstructor[*Person](b, func() *Person {
			return &Person{Name: "default"}
		})
		assert.NoError(t, err)

		codefixture.AddModel[*Person](b, func(p *Person) {
			assert.Equal(t, "default", p.Name)
		})
	})
}

func TestAddModel(t *testing.T) {
	t.Run("no setter, has constructor", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		err := codefixture.RegisterConstructor[*Person](b, func() *Person {
			return &Person{Name: "default"}
		})
		assert.NoError(t, err)

		ref, err := codefixture.AddModel[*Person](b, nil)
		assert.NoError(t, err)
		assert.Equal(t, "default", b.GetModel(ref).(*Person).Name)
	})
	t.Run("no setter, no constructor", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		ref, err := codefixture.AddModel[*Person](b, nil)
		assert.NoError(t, err)
		assert.Zero(t, b.GetModel(ref).(*Person).Name)
	})
	t.Run("override value by setter", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		err := codefixture.RegisterConstructor[*Person](b, func() *Person {
			return &Person{Name: "default"}
		})
		assert.NoError(t, err)

		ref, err := codefixture.AddModel[*Person](b, func(p *Person) {
			p.Name = "override"
		})
		assert.NoError(t, err)

		assert.Equal(t, "override", b.GetModel(ref).(*Person).Name)
	})
}
