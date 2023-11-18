package codefixture_test

import (
	"testing"

	"github.com/bluegreenhq/codefixture"
	"github.com/stretchr/testify/assert"
)

type PersonMaterial struct {
	Name    string
	GroupID int
}

func TestRegisterWriter(t *testing.T) {
	t.Run("in and out are same type", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		err := codefixture.RegisterWriter(b, func(p *Person) (*Person, error) {
			return p, nil
		})
		assert.NoError(t, err)
	})
	t.Run("in and out are different type", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		err := codefixture.RegisterWriter(b, func(pm *PersonMaterial) (*Person, error) {
			p := &Person{Name: pm.Name}
			return p, nil
		})
		assert.NoError(t, err)
	})
}

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

		m := codefixture.GetBuilderModel[*Person](b, ref)
		assert.Equal(t, "default", m.Name)
	})
	t.Run("no setter, no constructor", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		ref, err := codefixture.AddModel[*Person](b, nil)
		assert.NoError(t, err)

		m := codefixture.GetBuilderModel[*Person](b, ref)
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

		m := codefixture.GetBuilderModel[*Person](b, ref)
		assert.Equal(t, "override", m.Name)
	})
}

func TestConvertAndAddModel(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		err := codefixture.RegisterWriter(b, func(pm *PersonMaterial) (*Person, error) {
			p := &Person{Name: pm.Name}
			return p, nil
		})
		assert.NoError(t, err)

		ref, err := codefixture.ConvertAndAddModel[*PersonMaterial, *Person](b, func(p *PersonMaterial) {
			p.Name = "override"
		})
		assert.NoError(t, err)

		f, err := b.Build()
		assert.NoError(t, err)

		m := codefixture.GetModel[*Person](f, ref.ModelRef())
		assert.Equal(t, "override", m.Name)
	})
}

func TestAddModelWithRelation(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		b.RegisterWriter(&Person{}, func(p any) (any, error) {
			return p, nil
		})
		b.RegisterWriter(&Group{}, func(g any) (any, error) {
			return g, nil
		})
		err := codefixture.RegisterConstructor(b, func() *Person {
			return &Person{Name: "default"}
		})
		assert.NoError(t, err)

		g, err := codefixture.AddModel(b, func(g *Group) {
			g.ID = 1
			g.Name = "family"
		})
		assert.NoError(t, err)

		p, err := codefixture.AddModelWithRelation[*Person, *Group](b, g, func(p any, g any) {
			p.(*Person).GroupID = g.(*Group).ID
		})
		assert.NoError(t, err)

		f, err := b.Build()
		assert.NoError(t, err)

		m := codefixture.GetModel[*Person](f, p.ModelRef())
		assert.Equal(t, 1, m.GroupID)
	})
	t.Run("convert model", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		b.RegisterWriter(&PersonMaterial{}, func(p any) (any, error) {
			return &Person{}, nil
		})
		b.RegisterWriter(&Group{}, func(g any) (any, error) {
			return g, nil
		})

		p, err := codefixture.AddModel(b, func(p *PersonMaterial) {
			p.Name = "john"
		})
		assert.NoError(t, err)
		g, err := codefixture.AddModelWithRelation[*Group, *PersonMaterial](b, p, func(g any, p any) {
			g.(*Group).CreatorID = p.(*Person).ID
		})
		assert.NoError(t, err)
		group := codefixture.GetBuilderModel[*Group](b, g)
		assert.Equal(t, 0, group.CreatorID)
	})
}

func TestGetBuilderModel(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		b := codefixture.NewFixtureBuilder()
		err := codefixture.RegisterConstructor(b, func() *Person {
			return &Person{Name: "default"}
		})
		assert.NoError(t, err)

		ref, err := codefixture.AddModel[*Person](b, nil)
		assert.NoError(t, err)

		m := codefixture.GetBuilderModel(b, ref)
		assert.Equal(t, "default", m.Name)
	})
}
