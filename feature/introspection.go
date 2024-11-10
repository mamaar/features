package feature

type SchemaIntrospector struct {
	sch Schema
}

type IntrospectedField struct {
	exists bool
	field  Field
}

func NewSchemaIntrospector(sch Schema) *SchemaIntrospector {
	return &SchemaIntrospector{
		sch: sch,
	}
}

func (i *SchemaIntrospector) GetField(name string) (IntrospectedField, error) {
	var field Field
	for _, mig := range i.sch.Migrations {
		for _, op := range mig.Operations {
			if af, ok := op.(AddField); ok {
				if af.Field.Name == name {
					field = af.Field
				}
			}
		}
	}
	return IntrospectedField{
		exists: field.Name != "",
		field:  field,
	}, nil
}

func (f *IntrospectedField) Exists() bool {
	return f.exists
}

func (f *IntrospectedField) Type() FieldType {
	return f.field.Type
}
