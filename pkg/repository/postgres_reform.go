package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	reform "gopkg.in/reform.v1"

	"search-filter/pkg/models"
	"search-filter/pkg/types"
)

type PostgresRepository struct{ db *reform.DB }

func NewPostgresRepository(db *reform.DB) *PostgresRepository { return &PostgresRepository{db: db} }

func (r *PostgresRepository) Create(ctx context.Context, name string, query types.Query) (*models.Filter, error) {
	now := time.Now().UTC()
	f := &models.Filter{
		Name:      name,
		Query:     query,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := r.db.WithContext(ctx).Insert(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]models.FilterListItem, error) {
	rows, err := r.db.WithContext(ctx).FindAllFrom(models.FilterTable, "ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	res := make([]models.FilterListItem, 0, len(rows))
	for _, s := range rows {
		f := s.(*models.Filter)
		res = append(res, f.ToListItem())
	}
	return res, nil
}

func (r *PostgresRepository) Get(ctx context.Context, id uuid.UUID) (*models.Filter, error) {
	var f models.Filter
	if err := r.db.WithContext(ctx).FindByPrimaryKeyTo(&f, id); err != nil {
		if errors.Is(err, reform.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &f, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id uuid.UUID, query types.Query) (*models.Filter, error) {
	var f models.Filter
	if err := r.db.WithContext(ctx).FindByPrimaryKeyTo(&f, id); err != nil {
		if errors.Is(err, reform.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	f.Query = query // ← снова напрямую; Scanner/Valuer сделают работу
	f.UpdatedAt = time.Now().UTC()

	if err := r.db.WithContext(ctx).Update(&f); err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	n, err := r.db.WithContext(ctx).DeleteFrom(models.FilterTable, "WHERE id = $1", id)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
