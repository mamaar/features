package form

import (
	"net/url"
	"testing"

	"github.com/mamaar/jsonchamp"

	"github.com/mamaar/features/feature"
)

func newTestSchema(fields ...feature.Field) feature.Schema {
	ops := make([]feature.Operation, len(fields))
	for i, f := range fields {
		ops[i] = feature.AddField{Field: f}
	}
	return feature.Schema{
		Migrations: feature.Migrations{
			{Operations: ops},
		},
	}
}

func TestSetFromUrlValues_StringField(t *testing.T) {
	sch := newTestSchema(feature.Field{Name: "name", Type: feature.FieldTypeString})
	fe := feature.New(sch, feature.WithMap(jsonchamp.New()))
	fo := New(fe)

	err := fo.SetFromUrlValues(url.Values{"name": []string{"alice"}})
	if err != nil {
		t.Fatal(err)
	}

	got, err := fe.GetString("name")
	if err != nil {
		t.Fatal(err)
	}
	if got != "alice" {
		t.Fatalf("expected \"alice\", got %q", got)
	}
}

func TestSetFromUrlValues_NumberField(t *testing.T) {
	sch := newTestSchema(feature.Field{Name: "score", Type: feature.FieldTypeNumber})
	fe := feature.New(sch, feature.WithMap(jsonchamp.New()))
	fo := New(fe)

	err := fo.SetFromUrlValues(url.Values{"score": []string{"3.14"}})
	if err != nil {
		t.Fatal(err)
	}

	got, err := fe.GetFloat("score")
	if err != nil {
		t.Fatal(err)
	}
	if got != 3.14 {
		t.Fatalf("expected 3.14, got %v", got)
	}
}

func TestSetFromUrlValues_IntegerField(t *testing.T) {
	sch := newTestSchema(feature.Field{Name: "count", Type: feature.FieldTypeInteger})
	fe := feature.New(sch, feature.WithMap(jsonchamp.New()))
	fo := New(fe)

	err := fo.SetFromUrlValues(url.Values{"count": []string{"42"}})
	if err != nil {
		t.Fatal(err)
	}

	got, err := fe.GetInt("count")
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Fatalf("expected 42, got %v", got)
	}
}

func TestSetFromUrlValues_InvalidNumber(t *testing.T) {
	sch := newTestSchema(feature.Field{Name: "score", Type: feature.FieldTypeNumber})
	fe := feature.New(sch, feature.WithMap(jsonchamp.New()))
	fo := New(fe)

	err := fo.SetFromUrlValues(url.Values{"score": []string{"abc"}})
	if err == nil {
		t.Fatal("expected error for invalid number, got nil")
	}

	_, ok := fe.Get("score")
	if ok {
		t.Fatal("value should not be stored on parse failure")
	}
}

func TestSetFromUrlValues_InvalidInteger(t *testing.T) {
	sch := newTestSchema(feature.Field{Name: "count", Type: feature.FieldTypeInteger})
	fe := feature.New(sch, feature.WithMap(jsonchamp.New()))
	fo := New(fe)

	err := fo.SetFromUrlValues(url.Values{"count": []string{"nope"}})
	if err == nil {
		t.Fatal("expected error for invalid integer, got nil")
	}

	_, ok := fe.Get("count")
	if ok {
		t.Fatal("value should not be stored on parse failure")
	}
}

func TestSetFromUrlValues_UnknownFieldSkipped(t *testing.T) {
	sch := newTestSchema(feature.Field{Name: "name", Type: feature.FieldTypeString})
	fe := feature.New(sch, feature.WithMap(jsonchamp.New()))
	fo := New(fe)

	err := fo.SetFromUrlValues(url.Values{"unknown": []string{"val"}})
	if err != nil {
		t.Fatalf("unknown fields should not cause errors, got: %v", err)
	}

	_, ok := fe.Get("unknown")
	if ok {
		t.Fatal("unknown field should not be stored")
	}
}

func TestSetFromUrlValues_MultipleFields(t *testing.T) {
	sch := newTestSchema(
		feature.Field{Name: "name", Type: feature.FieldTypeString},
		feature.Field{Name: "score", Type: feature.FieldTypeNumber},
		feature.Field{Name: "count", Type: feature.FieldTypeInteger},
	)
	fe := feature.New(sch, feature.WithMap(jsonchamp.New()))
	fo := New(fe)

	err := fo.SetFromUrlValues(url.Values{
		"name":  []string{"bob"},
		"score": []string{"9.5"},
		"count": []string{"7"},
	})
	if err != nil {
		t.Fatal(err)
	}

	name, err := fe.GetString("name")
	if err != nil {
		t.Fatal(err)
	}
	if name != "bob" {
		t.Fatalf("expected \"bob\", got %q", name)
	}

	score, err := fe.GetFloat("score")
	if err != nil {
		t.Fatal(err)
	}
	if score != 9.5 {
		t.Fatalf("expected 9.5, got %v", score)
	}

	count, err := fe.GetInt("count")
	if err != nil {
		t.Fatal(err)
	}
	if count != 7 {
		t.Fatalf("expected 7, got %v", count)
	}
}

func TestFeatureGetter(t *testing.T) {
	sch := newTestSchema(feature.Field{Name: "x", Type: feature.FieldTypeString})
	fe := feature.New(sch, feature.WithMap(jsonchamp.New()))
	fo := New(fe)

	if fo.Feature() != fe {
		t.Fatal("Feature() should return the underlying feature")
	}
}
