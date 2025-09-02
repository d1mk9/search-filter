package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Query map[string]any

var emptyJSON = []byte("{}")

func (q Query) Value() (driver.Value, error) {
	if len(q) == 0 {
		return emptyJSON, nil
	}
	b, err := json.Marshal(q)
	if err != nil {
		return nil, fmt.Errorf("marshal Query: %w", err)
	}
	return b, nil
}
func (q *Query) Scan(src any) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("types.Query scan: unsupported source type")
	}
	return json.Unmarshal(b, q)
}
