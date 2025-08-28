package models

import (
	"encoding/json"
	"time"
)

//go:generate reform
//reform:filters
type Filter struct {
	ID        int64           `reform:"id,pk"      json:"id"`
	Name      string          `reform:"name"       json:"name"`
	Query     json.RawMessage `reform:"query"      json:"query"`
	CreatedAt time.Time       `reform:"created_at" json:"created_at"`
	UpdatedAt time.Time       `reform:"updated_at" json:"updated_at"`
}

type FilterListItem struct {
	ID    int64           `json:"id"`
	Name  string          `json:"name"`
	Query json.RawMessage `json:"query"`
}

func (f *Filter) ToListItem() FilterListItem {
	return FilterListItem{
		ID:    f.ID,
		Name:  f.Name,
		Query: f.Query,
	}
}
