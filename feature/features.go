package feature

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mamaar/jsonchamp/maps"
)

var (
	ErrPropertyNotFound = errors.New("property not found")
)

// Feature is a map of values with a specific schema.
type Feature struct {
	schema        Schema
	schemaVersion int
	m             *maps.Map
}

type Option func(*Feature)

func WithMap(m *maps.Map) Option {
	return func(f *Feature) {
		f.m = m
	}
}

func New(sch Schema, opts ...Option) *Feature {
	f := &Feature{
		schema: sch,
		m:      maps.New(),
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

func (f *Feature) Map() *maps.Map {
	return f.m
}

func (f *Feature) Schema() Schema {
	return f.schema
}

func (f *Feature) Set(key string, value any) {
	f.m = f.m.Set(key, value)
}

func (f *Feature) Get(key string) (any, bool) {
	return f.m.Get(key)
}

func (f *Feature) GetString(key string) (string, error) {
	v, err := f.m.GetString(key)
	if errors.Is(err, maps.ErrKeyNotFound) {
		return "", fmt.Errorf("%w: %s", ErrPropertyNotFound, key)
	}
	if err != nil {
		return "", err
	}
	return v, nil
}

func (f *Feature) GetFloat(key string) (float64, error) {
	v, err := f.m.GetFloat(key)
	if errors.Is(err, maps.ErrKeyNotFound) {
		return 0.0, fmt.Errorf("%w: %s", ErrPropertyNotFound, key)
	}
	if err != nil {
		return 0.0, err
	}
	return v, nil
}

func (f *Feature) GetInt(key string) (int, error) {
	v, err := f.m.GetInt(key)
	if errors.Is(err, maps.ErrKeyNotFound) {
		return 0, fmt.Errorf("%w: %s", ErrPropertyNotFound, key)
	}
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (f *Feature) SetInt(key string, value int) error {
	f.m = f.m.Set(key, value)
	return nil
}

func (f *Feature) Migrate(s Schema) error {
	var err error

	if f.schemaVersion >= len(s.Migrations) {
		return ErrSchemaVersionNotFound
	}
	for i, m := range s.Migrations[f.schemaVersion:] {
		for _, op := range m.Operations {
			f.m, err = op.Apply(f.m)
			if err != nil {
				return err
			}
		}
		f.schemaVersion = i + 1
	}
	return nil
}

func (f *Feature) UnmarshalJSON(data []byte) error {
	type decode struct {
		SchemaURN     string          `json:"schema_urn"`
		SchemaVersion int             `json:"schema_version"`
		Payload       json.RawMessage `json:"payload"`
	}
	var d decode
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}
	if d.Payload == nil {
		d.Payload = []byte("{}")
	}
	m := maps.New()
	if err := json.Unmarshal(d.Payload, &m); err != nil {
		return err
	}
	f.m = m
	return nil
}
