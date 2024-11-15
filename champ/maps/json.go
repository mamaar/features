package maps

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func marshalValue(m any) ([]byte, error) {
	buf := &bytes.Buffer{}
	switch v := m.(type) {
	case string:
		js, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("error marshalling string: %w", err)
		}
		buf.WriteString(string(js))
	case []string:
		js, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("error marshalling string array: %w", err)
		}
		buf.WriteString(string(js))
	case int:
		j, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("error marshalling int: %w", err)
		}
		buf.WriteString(string(j))
	case bool:
		j, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("error marshalling bool: %w", err)
		}
		buf.WriteString(string(j))
	case float64:
		j, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("error marshalling float64: %w", err)
		}
		buf.WriteString(string(j))
	case *Map:
		buf.WriteString("{")
		keys := v.Keys()
		for i, k := range keys {
			buf.WriteString(string(marshalKey(k)))
			buf.WriteString(":")
			v, _ := v.Get(k)
			d, err := marshalValue(v)
			if err != nil {
				return nil, err
			}
			buf.WriteString(string(d))
			if i < len(keys)-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString("}")
	case nil:
		buf.WriteString("null")
	default:
		return nil, fmt.Errorf("could not marshal type: %T", v)
	}
	return buf.Bytes(), nil
}

func marshalKey(m string) []byte {
	key, err := marshalValue(m)
	if err != nil {
		return nil
	}
	return key
}

func marshal(m *Map) ([]byte, error) {
	return marshalValue(m)
}

func (m *Map) MarshalJSON() ([]byte, error) {
	b, err := marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalArray(dec *json.Decoder, m *Map, key string) (*Map, error) {
	var arr []any
	for {
		token, err := dec.Token()
		if err != nil {
			return New(), fmt.Errorf("could not get token: %w", err)
		}
		if token == json.Delim(']') {
			return m.Set(key, arr), nil
		}
		switch v := token.(type) {
		case json.Delim:
			switch v {
			case '{':
				newMap := New()
				newMap, err := unmarshalMap(dec, newMap)
				if err != nil {
					return New(), err
				}
				arr = append(arr, newMap)
			case '}':
				return m.Set(key, arr), nil
			default:
				return New(), fmt.Errorf("unexpected delimiter %c", v)
			}
		default:
			return New(), fmt.Errorf("unexpected type %T", v)
		}
	}
}

func unmarshalMap(dec *json.Decoder, m *Map) (*Map, error) {
	for {
		keyToken, err := dec.Token()
		if err != nil {
			return nil, fmt.Errorf("could not get key: %w", err)
		}
		if keyToken == json.Delim('}') {
			return m, nil
		}
		keyString, isString := keyToken.(string)
		if !isString {
			return nil, fmt.Errorf("expected string key, got %T", keyToken)
		}
		valueToken, err := dec.Token()
		if err != nil {
			return nil, fmt.Errorf("could not get value: %w", err)
		}
		switch v := valueToken.(type) {
		case string:
			m = m.Set(keyString, v)
		case float64:
			m = m.Set(keyString, v)
		case bool:
			m = m.Set(keyString, v)
		case json.Delim:
			switch v {
			case '{':
				newMap := New()
				newMap, err := unmarshalMap(dec, newMap)
				if err != nil {
					return nil, err
				}
				m = m.Set(keyString, newMap)
			case '[':
				return unmarshalArray(dec, m, keyString)
			default:
				return nil, fmt.Errorf("unexpected delimiter %c", v)
			}
		default:
			return nil, fmt.Errorf("unexpected type %T", v)
		}
	}
}

func (m *Map) UnmarshalJSON(d []byte) error {
	if m.hasher == nil {
		m.hasher = DefaultMapOptions.hasher()
	}
	if m.root == nil {
		m.root = &bitmasked{
			level:      0,
			valueMap:   0,
			subMapsMap: 0,
			values:     []node{},
		}
	}
	dec := json.NewDecoder(bytes.NewReader(d))
	token, err := dec.Token()
	if err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	if token != json.Delim('{') {
		return fmt.Errorf("expected object, got %T", token)
	}

	newMap, err := unmarshalMap(dec, m)
	if err != nil {
		return err
	}
	*m = *newMap
	return nil
}
