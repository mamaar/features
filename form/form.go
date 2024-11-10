package form

import (
	"net/url"
	"strconv"

	"github.com/mamaar/features/feature"
)

// New creates a new form with the given feature.

func New(fe *feature.Feature) *Form {

	return &Form{
		feat: fe,
	}

}

// Form represents a form with a feature.

type Form struct {
	feat *feature.Feature
}

// SetFromUrlValues sets the form values from the given url.Values.
func (f *Form) SetFromUrlValues(values url.Values) {
	schemaIntro := feature.NewSchemaIntrospector(f.feat.Schema())

	for key := range values {
		field, err := schemaIntro.GetField(key)
		if err != nil {
			continue
		}

		if !field.Exists() {
			continue
		}

		switch field.Type() {
		case feature.FieldTypeString:
			val := values.Get(key)
			f.feat.Set(key, val)
		case feature.FieldTypeNumber:
			flt, err := strconv.ParseFloat(values.Get(key), 64)
			if err != nil {
				f.feat.Set(key, values.Get(key))
			} else {
				f.feat.Set(key, flt)
			}
		}
	}

}

func (f *Form) Validate() error {
	compiled, err := f.feat.Schema().ToJSONSchema()
	if err != nil {
		return err
	}

	return compiled.Validate(f.feat)
}
