package feature

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/santhosh-tekuri/jsonschema/v6"

	"github.com/mamaar/jsonchamp/maps"
)

var (
	ErrInvalidSchema         = errors.New("invalid schema")
	ErrSchemaVersionNotFound = errors.New("schema version not found")
)

type Schema struct {
	Schema     string     `json:"schema"`
	Migrations Migrations `json:"migrations"`
}

var Empty = Schema{
	Schema:     "",
	Migrations: Migrations{},
}

// Migrate applies the migrations to the input map.
func (s Schema) Migrate(m *maps.Map) (*maps.Map, error) {
	var err error

	res := m.Copy()
	for _, migration := range s.Migrations {
		for _, op := range migration.Operations {
			res, err = op.Apply(res)
			if err != nil {
				return nil, err
			}
		}
	}
	return res, nil
}

func (s Schema) ToJSONSchema() (*Validator, error) {
	schemaMap, err := s.Migrations.Reduce()
	if err != nil {
		return nil, err
	}
	js, err := json.Marshal(schemaMap)
	if err != nil {
		return nil, err
	}
	sch, err := jsonschema.UnmarshalJSON(bytes.NewReader(js))
	if err != nil {
		return nil, err
	}

	c := jsonschema.NewCompiler()
	err = c.AddResource("schema", sch)
	if err != nil {
		return nil, err
	}

	compiled, err := c.Compile("schema")
	if err != nil {
		return nil, err
	}

	return NewValidator(compiled), nil
}
