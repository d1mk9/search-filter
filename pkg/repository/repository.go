package repository

import (
	"context"
	"errors"

	"search-filter/pkg/models"
	"search-filter/pkg/types"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("filter not found")

type Repository interface {
	Create(ctx context.Context, name string, query types.Query) (*models.Filter, error)
	List(ctx context.Context) ([]models.FilterListItem, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Filter, error)
	Update(ctx context.Context, id uuid.UUID, query types.Query) (*models.Filter, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
