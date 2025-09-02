package models

import (
	"search-filter/pkg/types"

	"time"

	"github.com/google/uuid"
)

//go:generate reform
//reform:filters
type Filter struct {
	ID        uuid.UUID   `reform:"id,pk"      json:"id"`
	Name      string      `reform:"name"       json:"name"`
	Query     types.Query `reform:"query"      json:"query"`
	CreatedAt time.Time   `reform:"created_at" json:"created_at"`
	UpdatedAt time.Time   `reform:"updated_at" json:"updated_at"`
}

type FilterListItem struct {
	ID    uuid.UUID   `json:"id"`
	Name  string      `json:"name"`
	Query types.Query `json:"query"`
}

func (f *Filter) ToListItem() FilterListItem {
	return FilterListItem{
		ID:    f.ID,
		Name:  f.Name,
		Query: f.Query,
	}
}
