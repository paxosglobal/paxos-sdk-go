package paxos

import (
	"bytes"
	"encoding/json"
	"reflect"
)

type rawJSONSetter interface {
	setRawJSON(string)
}

func rawJSONBytes(s string) json.RawMessage {
	if s == "" {
		return nil
	}
	return json.RawMessage(s)
}

func attachRawJSON(v any, raw []byte) {
	if s, ok := v.(rawJSONSetter); ok {
		s.setRawJSON(string(bytes.TrimSpace(raw)))
	}
}

func decodeItems[T any](array json.RawMessage) ([]T, error) {
	dec := json.NewDecoder(bytes.NewReader(array))
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if tok == nil {
		return nil, nil
	}
	if delim, ok := tok.(json.Delim); !ok || delim != '[' {
		return nil, &json.UnmarshalTypeError{Value: tokenKind(tok), Type: reflect.TypeFor[[]T]()}
	}
	items := []T{}
	prev := int(dec.InputOffset())
	for dec.More() {
		start := prev
		for start < len(array) && (array[start] == ',' || array[start] == ' ' || array[start] == '\t' || array[start] == '\n' || array[start] == '\r') {
			start++
		}
		if start >= len(array) || array[start] != '{' {
			return nil, &json.UnmarshalTypeError{Value: "non-object element", Type: reflect.TypeFor[T]()}
		}
		items = append(items, *new(T))
		item := &items[len(items)-1]
		if err := dec.Decode(item); err != nil {
			return nil, err
		}
		end := int(dec.InputOffset())
		attachRawJSON(item, array[start:end])
		prev = end
	}
	return items, nil
}

func tokenKind(tok json.Token) string {
	switch tok.(type) {
	case string:
		return "string"
	case float64, json.Number:
		return "number"
	case bool:
		return "bool"
	}
	return "object"
}
