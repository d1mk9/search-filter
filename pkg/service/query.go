package service

import (
	"encoding/json"
	"fmt"
)

type Query map[string]any

func validateQueryObjectMap(q Query) error {
	if q == nil {
		return fmt.Errorf("query is required")
	}
	if len(q) == 0 {
		return fmt.Errorf("query must not be an empty object")
	}
	return nil
}

func marshalQuery(q Query) (json.RawMessage, error) {
	b, err := json.Marshal(q)
	if err != nil {
		return nil, fmt.Errorf("encode query to JSON: %w", err)
	}
	return json.RawMessage(b), nil
}

func unmarshalQuery(raw json.RawMessage) (Query, error) {
	var q Query
	if err := json.Unmarshal(raw, &q); err != nil {
		return nil, fmt.Errorf("decode JSON to query object: %w", err)
	}
	return q, nil
}
