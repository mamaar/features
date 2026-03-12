package form

import (
	"errors"
	"fmt"
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

// Feature returns the underlying feature.
func (f *Form) Feature() *feature.Feature {
	return f.feat
}

// SetFromUrlValues sets the form values from the given url.Values.
func (f *Form) SetFromUrlValues(values url.Values) error {
	schemaIntro := feature.NewSchemaIntrospector(f.feat.Schema())

	var errs []error
	for key := range values {
		field, err := schemaIntro.GetField(key)
		if err != nil {
			continue
		}

		if !field.Exists() {
			continue
		}

		raw := values.Get(key)
		switch field.Type() {
		case feature.FieldTypeString:
			f.feat.Set(key, raw)
		case feature.FieldTypeNumber:
			flt, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				errs = append(errs, fmt.Errorf("field %q: %w", key, err))
			} else {
				f.feat.Set(key, flt)
			}
		case feature.FieldTypeInteger:
			parsed, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				errs = append(errs, fmt.Errorf("field %q: %w", key, err))
			} else {
				f.feat.SetInt(key, int(parsed))
			}
		}
	}

	return errors.Join(errs...)
}

func (f *Form) Validate() error {
	compiled, err := f.feat.Schema().ToJSONSchema()
	if err != nil {
		return err
	}

	return compiled.Validate(f.feat)
}
