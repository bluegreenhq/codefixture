package codefixture

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixtureBuilder_getModelsOrderedByRelations(t *testing.T) {
	t.Run("no relations", func(t *testing.T) {
		b := NewFixtureBuilder()
		b.models["a"] = "a"
		b.models["b"] = "b"

		refs, ms := b.getModelsOrderedByRelations()
		assert.Len(t, refs, 2)
		assert.Len(t, ms, 2)
	})
	t.Run("leaf should appear first", func(t *testing.T) {
		b := NewFixtureBuilder()
		b.models["root"] = "root"
		b.models["branch"] = "branch"
		b.models["leaf"] = "leaf"
		b.relations = []ModelRelation{
			{
				TargetRef:  "root",
				ForeignRef: "branch",
			},
			{
				TargetRef:  "branch",
				ForeignRef: "leaf",
			},
		}

		refs, ms := b.getModelsOrderedByRelations()
		assert.Len(t, refs, 3)
		assert.Len(t, ms, 3)
		assert.Equal(t, "leaf", ms[0])
		assert.Equal(t, "branch", ms[1])
		assert.Equal(t, "root", ms[2])
	})
	t.Run("delayed dependency resolution", func(t *testing.T) {
		b := NewFixtureBuilder()
		b.models["d"] = "d"
		b.models["c"] = "c"
		b.models["b"] = "b"
		b.models["a"] = "a"

		b.relations = []ModelRelation{
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

		refs, ms := b.getModelsOrderedByRelations()
		assert.Len(t, refs, 4)
		assert.Len(t, ms, 4)
		assert.Equal(t, []any{"a", "b", "c", "d"}, ms)
	})
}
