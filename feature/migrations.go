package feature

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mamaar/features/champ/maps"
)

var (
	ErrCyclicDependency    = errors.New("cyclic dependency")
	ErrDependencyNotExists = errors.New("dependency does not exist")
	ErrInvalidMigration    = errors.New("invalid migration")
)

type Operation interface {
	// Apply applies the operation to the given data model.
	Apply(*maps.Map) (*maps.Map, error)
}

type FieldType string

const (
	FieldTypeString FieldType = "string"
	FieldTypeNumber FieldType = "number"
)

type Field struct {
	Name     string
	Type     FieldType
	Required bool
	Default  any
}

type AddField struct {
	Field Field
}

func NewAddFieldFromMap(m *maps.Map) (AddField, error) {
	fieldDef, err := m.GetMap("field")
	if err != nil {
		return AddField{}, err
	}

	name, err := fieldDef.GetString("name")
	if err != nil {
		return AddField{}, err
	}
	typ, err := fieldDef.GetString("type")
	if err != nil {
		return AddField{}, err
	}
	required, err := fieldDef.GetBool("required")
	if err != nil {
		return AddField{}, err
	}
	def, ok := fieldDef.Get("default")
	if !ok {
		def = nil
	}
	return AddField{
		Field: Field{
			Name:     name,
			Type:     FieldType(typ),
			Required: required,
			Default:  def,
		},
	}, nil
}

// Apply implements Operation.
// It handles migrations of the data model by adding a new field.
// If the field is required, it must have a default value.
func (a AddField) Apply(in *maps.Map) (*maps.Map, error) {
	hasCurrentValue := in.Contains(a.Field.Name)

	// If the field is required it must have a default value or already exist in the data model.
	if a.Field.Required && a.Field.Default == nil && !hasCurrentValue {
		return nil, fmt.Errorf("required field '%s' must have a default value", a.Field.Name)
	}

	// Set the default value if the field does not exist in the data model.
	if !hasCurrentValue && a.Field.Default != nil {
		return in.Set(a.Field.Name, a.Field.Default), nil
	}
	return in, nil
}

var _ Operation = AddField{}

type AlterField struct {
	Field Field
}

// Apply implements Operation.
func (a AlterField) Apply(*maps.Map) (*maps.Map, error) {
	panic("unimplemented")
}

var _ Operation = AlterField{}

type RemoveField struct {
	FieldName string
}

// Apply implements Operation.
func (r RemoveField) Apply(in *maps.Map) (*maps.Map, error) {
	n, wasDeleted := in.Delete(r.FieldName)
	if !wasDeleted {
		return nil, fmt.Errorf("field '%s' does not exist", r.FieldName)
	}
	return n, nil
}

var _ Operation = RemoveField{}

type Migration struct {
	Description string      `json:"description"`
	Operations  []Operation `json:"operations"`
}

func unmarshalOperations(ops []any) ([]Operation, error) {
	opsRes := make([]Operation, len(ops))

	for i, op := range ops {
		opMap, ok := op.(*maps.Map)
		if !ok {
			return nil, fmt.Errorf("operation %d is not a map", i)
		}

		opType, ok := opMap.Get("type")
		if !ok {
			return nil, fmt.Errorf("operation %d does not have a type", i)
		}
		switch opType {
		case "add_field":
			addField, err := NewAddFieldFromMap(opMap)
			if err != nil {
				return nil, err
			}
			opsRes[i] = addField
		case "remove_field":
			var r RemoveField
			opsRes[i] = r
		default:
			return nil, fmt.Errorf("operation not implemented: %s", opType)
		}
	}

	return opsRes, nil
}

func (m *Migration) UnmarshalJSON(data []byte) error {
	var d *maps.Map
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	description, err := d.GetString("description")
	if err != nil {
		return err
	}

	operations, ok := d.Get("operations")
	if !ok {
		return fmt.Errorf("migration does not have operations")
	}

	ops, err := unmarshalOperations(operations.([]any))
	if err != nil {
		return err
	}

	m.Description = description
	m.Operations = ops

	return nil
}

type Migrations []*Migration

func (m Migrations) Validate() error {
	return nil
}

// Reduce returns a JSON schema as a map that can be used to validate the data.
func (m Migrations) Reduce() (*maps.Map, error) {
	required := maps.New()

	properties := maps.NewFromItems()
	for migrationIndex, migration := range m {
		for _, op := range migration.Operations {
			switch op := op.(type) {
			case AddField:
				field := op.Field
				if len(field.Name) == 0 {
					return nil, errors.New("field name must not be empty")
				}
				if fieldExists := properties.Contains(field.Name); fieldExists {
					return nil, fmt.Errorf("field '%s' already exists", field.Name)
				}
				// Required fields must have a default value, unless it's the first migration.
				if (op.Field.Required && op.Field.Default == nil) && (migrationIndex != 0) {
					return nil, fmt.Errorf("required field must have a default value: %s", op.Field.Name)
				}

				properties = properties.Set(field.Name, maps.NewFromItems("type", string(field.Type)))
				if op.Field.Required {
					required = required.Set(field.Name, true)
				}

			case RemoveField:
				field := op.FieldName
				var propertyWasDeleted bool
				properties, propertyWasDeleted = properties.Delete(field)
				if !propertyWasDeleted {
					return nil, fmt.Errorf("field '%s' does not exist", field)
				}
				required, _ = required.Delete(field)

			default:
				return nil, fmt.Errorf("operation not implemented: %T", op)
			}
		}
	}

	schemaDocument := maps.NewFromItems(
		"properties", properties,
		"required", required.Keys(),
	)

	return schemaDocument, nil
}
