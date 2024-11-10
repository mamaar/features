package feature

import "fmt"

type SchemaStore interface {
	Get(schemaUrn string) (Schema, error)
}

type SchemaStoreMock struct {
	GetFn func(schemaUrn string) (Schema, error)
}

func (s SchemaStoreMock) Get(schemaUrn string) (Schema, error) {
	if s.GetFn != nil {
		return s.GetFn(schemaUrn)
	}
	return Schema{}, fmt.Errorf("GetFn is not implemented")
}

// Store is a feature storage interface
type Store interface {
	Get(key string) (*Feature, error)
}

// StoreMock is a mock implementation of the Store interface
type StoreMock struct {
	GetFn func(key string) (*Feature, error)
}

func (s *StoreMock) Get(key string) (*Feature, error) {
	if s.GetFn != nil {
		return s.GetFn(key)
	}
	return nil, fmt.Errorf("GetFn is not implemented")
}
