package reconcile

import (
	"testing"

	"github.com/mamaar/jsonchamp"
)

func TestReconcile(t *testing.T) {
	tests := []struct {
		name     string
		incoming *jsonchamp.Map
		base     *jsonchamp.Map
		head     *jsonchamp.Map
	}{
		{
			name:     "nothing changed",
			incoming: jsonchamp.NewFromItems("a", 1),
			base:     jsonchamp.NewFromItems("a", 1),
			head:     jsonchamp.NewFromItems("a", 1),
		},
		{
			name:     "add new field",
			incoming: jsonchamp.NewFromItems("a", 1, "b", 2),
			base:     jsonchamp.NewFromItems("a", 1),
			head:     jsonchamp.NewFromItems("a", 1),
		},
		{
			name:     "modify existing field",
			incoming: jsonchamp.NewFromItems("a", 2),
			base:     jsonchamp.NewFromItems("a", 1),
			head:     jsonchamp.NewFromItems("a", 1),
		},
		{
			name:     "delete field",
			incoming: jsonchamp.NewFromItems(),
			base:     jsonchamp.NewFromItems("a", 1),
			head:     jsonchamp.NewFromItems("a", 1),
		},
		{
			name:     "incoming matches head but not base",
			incoming: jsonchamp.NewFromItems("a", 2),
			base:     jsonchamp.NewFromItems("a", 1),
			head:     jsonchamp.NewFromItems("a", 2),
		},
		{
			name:     "incoming and head have different changes",
			incoming: jsonchamp.NewFromItems("a", 2),
			base:     jsonchamp.NewFromItems("a", 1),
			head:     jsonchamp.NewFromItems("a", 3),
		},
		{
			name:     "nested object changes",
			incoming: jsonchamp.NewFromItems("a", jsonchamp.NewFromItems("b", 2)),
			base:     jsonchamp.NewFromItems("a", jsonchamp.NewFromItems("b", 1)),
			head:     jsonchamp.NewFromItems("a", jsonchamp.NewFromItems("b", 1)),
		},
		{
			name:     "array changes",
			incoming: jsonchamp.NewFromItems("a", []interface{}{1, 2, 3}),
			base:     jsonchamp.NewFromItems("a", []interface{}{1, 2}),
			head:     jsonchamp.NewFromItems("a", []interface{}{1, 2}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reconcile(tt.incoming, tt.base, tt.head)
		})
	}
}
