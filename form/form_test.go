package form

import (
	"net/url"
	"testing"

	"github.com/mamaar/features/champ/maps"
	"github.com/mamaar/features/feature"
)

func TestSomething(t *testing.T) {

	sch := feature.Schema{
		Migrations: feature.Migrations{
			{
				Operations: []feature.Operation{
					feature.AddField{
						Field: feature.Field{
							Name:     "name",
							Type:     feature.FieldTypeNumber,
							Required: false,
						},
					},
				},
			},
		},
	}

	fe := feature.New(
		sch,
		feature.WithMap(maps.New()),
	)

	fo := New(fe)

	err := fo.Validate()
	if err != nil {
		t.Fatal(err)
	}

	values := url.Values{
		"name": []string{"1"},
	}

	fo.SetFromUrlValues(values)

	err = fo.Validate()
	if err != nil {
		t.Fatal(err)
	}

	val, err := fo.feat.GetInt("name")
	if err != nil {
		t.Fatal(err)
	}
	if val != 1 {
		t.Fatalf("expected 1, got %v", val)
	}

}
