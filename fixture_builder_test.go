package codefixture

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
}
