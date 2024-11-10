package database

import (
	"strings"

	"github.com/mamaar/features/champ/maps"
	"github.com/mamaar/features/feature"
)

const (
	CompositeSeparator = ":"
)

type IndexMap struct {
}

type Identifier string

type IdentifierBuilder func(f *feature.Feature) (Identifier, error)

func Literal(id Identifier) IdentifierBuilder {
	return func(f *feature.Feature) (Identifier, error) {
		return id, nil
	}
}

func Prop(propName string) IdentifierBuilder {
	return func(f *feature.Feature) (Identifier, error) {
		value, err := f.GetString(propName)
		if err != nil {
			return "", err
		}
		return Identifier(value), nil
	}
}

func Comp(items ...IdentifierBuilder) IdentifierBuilder {
	return func(f *feature.Feature) (Identifier, error) {
		identifiers := make([]string, 0, len(items))
		for _, item := range items {
			identifier, err := item(f)
			if err != nil {
				return "", err
			}
			identifiers = append(identifiers, string(identifier))
		}
		return Identifier(strings.Join(identifiers, CompositeSeparator)), nil
	}
}

type IndexTuple struct {
	Pk Identifier
	Sk Identifier
}

func CreateIndexedDatabasePayload(feat *feature.Feature, indexMap map[string]IdentifierBuilder) (*maps.Map, error) {
	payload := maps.New()
	for col, index := range indexMap {
		identifier, err := index(feat)
		if err != nil {
			return nil, err
		}
		payload = payload.Set(col, string(identifier))
	}
	payload = payload.Set("payload", feat.Map())
	return payload, nil
}
