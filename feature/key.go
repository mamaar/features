package feature

import (
	"strings"

	"github.com/mamaar/jsonchamp"
)

const (
	KeySeparator = ":"
)

type Key string

type KeyFunc func(f *Feature) (Key, error)

func LiteralKey(id Key) KeyFunc {
	return func(f *Feature) (Key, error) {
		return id, nil
	}
}

func PropKey(propName string) KeyFunc {
	return func(f *Feature) (Key, error) {
		value, err := f.GetString(propName)
		if err != nil {
			return "", err
		}
		return Key(value), nil
	}
}

func CompositeKey(items ...KeyFunc) KeyFunc {
	return func(f *Feature) (Key, error) {
		identifiers := make([]string, 0, len(items))
		for _, item := range items {
			identifier, err := item(f)
			if err != nil {
				return "", err
			}
			identifiers = append(identifiers, string(identifier))
		}
		return Key(strings.Join(identifiers, KeySeparator)), nil
	}
}

func CreateKeyedPayload(feat *Feature, indexMap map[string]KeyFunc) (*jsonchamp.Map, error) {
	payload := jsonchamp.New()
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
