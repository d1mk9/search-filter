package repository

import (
	"context"
	"encoding/json"
	"errors"

	"search-filter/pkg/models"
)

var ErrNotFound = errors.New("filter not found")

type Repository interface {
	Create(ctx context.Context, name string, query json.RawMessage) (*models.Filter, error)
	List(ctx context.Context) ([]models.FilterListItem, error)
	Get(ctx context.Context, id int64) (*models.Filter, error)
	Update(ctx context.Context, id int64, query json.RawMessage) (*models.Filter, error)
	Delete(ctx context.Context, id int64) error
}
