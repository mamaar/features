package feature

import "testing"

func TestGetField(t *testing.T) {

	sch := Schema{
		Migrations: Migrations{
			{
				Operations: []Operation{
					AddField{
						Field: Field{
							Name: "name",
							Type: FieldTypeString,
						},
					},
				},
			},
		},
	}

	intro := NewSchemaIntrospector(sch)

	field, err := intro.GetField("name")
	if err != nil {
		t.Fatalf("GetField(%q) = %v; want nil", "name", err)
	}

	if !field.Exists() {
		t.Fatalf("GetField(%q).Exists() = false; want true", "name")
	}

	typ := field.Type()
	if typ != FieldTypeString {
		t.Errorf("GetField(%q).Type() = %v; want %v", "name", typ, FieldTypeString)
	}

}
