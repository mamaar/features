package feature

import (
	"testing"
)

func TestBuildsSchema(t *testing.T) {

	migrations := Migrations{
		&Migration{
			Operations: []Operation{
				AddField{
					Field: Field{Name: "family", Type: "enum", Required: true, Default: "cat"},
				},
			},
		},
	}

	_, err := migrations.Reduce()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveFieldMigration(t *testing.T) {

	migrations := Migrations{
		&Migration{
			Operations: []Operation{
				AddField{
					Field: Field{
						Name:     "family",
						Type:     FieldTypeString,
						Required: false,
						Default:  "",
					},
				},
				RemoveField{
					FieldName: "family",
				},
			},
		},
	}

	_, err := migrations.Reduce()
	if err != nil {
		t.Fatal(err)
	}
}
