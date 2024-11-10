package feature

import (
	"encoding/json"
	"errors"
	"testing"
)

var orderSchemaMigrations = `{
	"schema": "urn:features:order",
	"migrations": [
		{
			"description": "Initial schema",
			"operations": [
				{
					"type": "add_field",
					"field": {
						"name": "order_id",
						"type": "number",
						"required": true
					}
				}
			]
		},
		{
			"description": "Add customer_id",
			"operations": [
				{
					"type": "add_field",
					"field": {	
						"name": "customer_id",	
						"type": "number",	
						"required": true,
						"default": 0
					}
				}
			]
		}
	]
}`

var order1Data = `
{
	"schema_urn": "urn:features:order/1",
	"payload": {
		"order_id": 123.456
	}
}`

func TestFeatures(t *testing.T) {
	schemaStore := &SchemaStoreMock{
		GetFn: func(schemaUrn string) (Schema, error) {
			switch schemaUrn {
			case "urn:features:order":
				var sch Schema
				err := json.Unmarshal([]byte(orderSchemaMigrations), &sch)
				if err != nil {
					return Empty, err
				}
				return sch, nil
			}
			return Empty, errors.New("schema not found")
		},
	}

	featureStore := &StoreMock{
		GetFn: func(key string) (*Feature, error) {
			switch key {
			case "urn:features:order/1":
				var feat *Feature
				err := json.Unmarshal([]byte(order1Data), &feat)
				if err != nil {
					return nil, err
				}
				return feat, nil
			}
			return nil, errors.New("feature not found")
		},
	}

	sch, err := schemaStore.Get("urn:features:order")
	if err != nil {
		t.Fatal(err)
	}

	feat, err := featureStore.Get("urn:features:order/1")
	if err != nil {
		t.Fatal(err)
	}

	// We have the latest schema and ensure the feature is migrated to it.
	err = feat.Migrate(sch)
	if err != nil {
		t.Fatal(err)
	}

	vI, err := feat.GetInt("customer_id")
	if errors.Is(err, ErrPropertyNotFound) {
		t.Logf("customer ID does not exist on feature, so we do something to handle that case...")
	}
	t.Logf("customer_id: %v, err: %v", vI, err)
	vS, ok := feat.GetFloat("order_id")
	t.Logf("order_id: %v, ok: %v", vS, ok)

	err = feat.SetInt("customer_id", 123)
	if err != nil {
		t.Fatal(err)
	}

	validator, err := sch.ToJSONSchema()
	if err != nil {
		t.Fatal(err)
	}

	err = validator.Validate(feat)
	if err != nil {
		t.Fatal(err)
	}
}
