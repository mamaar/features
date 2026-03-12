package reconcile

import (
	"testing"

	"github.com/mamaar/features/feature"
	"github.com/mamaar/jsonchamp"
)

func feat(kvs ...interface{}) *feature.Feature {
	return feature.New(feature.Schema{}, feature.WithMap(jsonchamp.NewFromItems(kvs...)))
}

func TestReconcile(t *testing.T) {
	tests := []struct {
		name     string
		incoming *feature.Feature
		base     *feature.Feature
		head     *feature.Feature
	}{
		{
			name:     "nothing changed",
			incoming: feat("a", 1),
			base:     feat("a", 1),
			head:     feat("a", 1),
		},
		{
			name:     "add new field",
			incoming: feat("a", 1, "b", 2),
			base:     feat("a", 1),
			head:     feat("a", 1),
		},
		{
			name:     "modify existing field",
			incoming: feat("a", 2),
			base:     feat("a", 1),
			head:     feat("a", 1),
		},
		{
			name:     "delete field",
			incoming: feat(),
			base:     feat("a", 1),
			head:     feat("a", 1),
		},
		{
			name:     "incoming matches head but not base",
			incoming: feat("a", 2),
			base:     feat("a", 1),
			head:     feat("a", 2),
		},
		{
			name:     "incoming and head have different changes",
			incoming: feat("a", 2),
			base:     feat("a", 1),
			head:     feat("a", 3),
		},
		{
			name:     "nested object changes",
			incoming: feat("a", jsonchamp.NewFromItems("b", 2)),
			base:     feat("a", jsonchamp.NewFromItems("b", 1)),
			head:     feat("a", jsonchamp.NewFromItems("b", 1)),
		},
		{
			name:     "array changes",
			incoming: feat("a", []interface{}{1, 2, 3}),
			base:     feat("a", []interface{}{1, 2}),
			head:     feat("a", []interface{}{1, 2}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reconcile(tt.incoming, tt.base, tt.head)
		})
	}
}
