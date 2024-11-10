package feature

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"

	"github.com/mamaar/features/feature/meta"
)

func TestSchemaSchema(t *testing.T) {
	tests := []struct {
		name       string
		migrations Migrations
	}{
		{
			name: "TestBuildsSchema",
			migrations: Migrations{
				&Migration{
					Operations: []Operation{
						AddField{
							Field: Field{Name: "family", Type: "enum", Required: true, Default: "cat"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := jsonschema.NewCompiler()

			dirEntries, err := meta.FS.ReadDir(".")
			if err != nil {
				t.Fatal(err)
			}

			for _, dirEntry := range dirEntries {
				if dirEntry.IsDir() {
					continue
				}
				name := dirEntry.Name()
				fileContents, err := meta.FS.ReadFile(name)
				if err != nil {
					t.Fatal(err)
				}

				fileData, err := jsonschema.UnmarshalJSON(bytes.NewReader(fileContents))
				if err != nil {
					t.Fatal(err)
				}

				err = compiler.AddResource(name, fileData)
				if err != nil {
					t.Fatal(err)
				}
			}

			metaSch, err := compiler.Compile("schema.json")
			if err != nil {
				t.Fatal(err)
			}

			schema, err := tt.migrations.Reduce()
			if err != nil {
				t.Fatal(err)
			}

			if schema == nil {
				t.Fatal("schema is nil")
			}

			schemaJSON, err := json.Marshal(schema)
			if err != nil {
				t.Fatal(err)
			}

			in, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaJSON))
			if err != nil {
				t.Fatal(err)
			}

			err = metaSch.Validate(in)
			if err != nil {
				t.Fatal(err)
			}

		})
	}
}
