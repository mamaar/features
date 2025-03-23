package feature

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"

	"github.com/mamaar/jsonchamp/maps"
)

type Validator struct {
	c *jsonschema.Schema
}

func NewValidator(s *jsonschema.Schema) *Validator {
	return &Validator{
		c: s,
	}
}

func (v *Validator) Validate(m *Feature) error {
	js, err := json.Marshal(m.m)
	if err != nil {
		return err
	}
	feat, err := jsonschema.UnmarshalJSON(strings.NewReader(string(js)))
	if err != nil {
		return err
	}
	err = v.c.Validate(feat)
	if err == nil {
		return nil
	}
	var validationErr *jsonschema.ValidationError
	if errors.As(err, &validationErr) {
		validationErrors := maps.New()
		for _, cause := range validationErr.Causes {
			switch k := cause.ErrorKind.(type) {
			case *kind.Type:
				location := cause.InstanceLocation[0]
				validationErrors = validationErrors.Set(location, maps.NewFromItems("got", k.Got, "expect", k.Want))
			case *kind.Required:
				for _, missing := range k.Missing {
					validationErrors = validationErrors.Set(missing, maps.NewFromItems("required", true))
				}
			default:
				panic("unexpected error kind")
			}
		}
		if len(validationErrors.Keys()) > 0 {
			return validationErrors
		}
	}
	if err != nil {
		return err
	}
	return nil
}
