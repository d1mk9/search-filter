package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"search-filter/pkg/models"
	"search-filter/pkg/placeholder"
	"search-filter/pkg/repository"
	"search-filter/pkg/types"

	"github.com/google/uuid"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrValidation = errors.New("validation error")
)

type Filters interface {
	Create(ctx context.Context, name string, query types.Query) (*models.Filter, error)
	List(ctx context.Context) ([]models.FilterListItem, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Filter, error)
	Update(ctx context.Context, id uuid.UUID, query types.Query) (*models.Filter, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Apply(ctx context.Context, id uuid.UUID) (types.Query, error)
}

type service struct {
	repo          repository.Repository
	loc           *time.Location
	currentUserID int64
}

func NewFiltersService(repo repository.Repository, loc *time.Location, currentUserID int64) (Filters, error) {
	if repo == nil {
		return nil, fmt.Errorf("NewFiltersService: repo is nil")
	}
	if loc == nil {
		return nil, fmt.Errorf("NewFiltersService: loc is nil")
	}
	if currentUserID <= 0 {
		return nil, fmt.Errorf("NewFiltersService: currentUserID must be > 0, got %d", currentUserID)
	}
	return &service{repo: repo, loc: loc, currentUserID: currentUserID}, nil
}

func (s *service) Create(ctx context.Context, name string, query types.Query) (*models.Filter, error) {
	return s.repo.Create(ctx, name, query)
}

func (s *service) List(ctx context.Context) ([]models.FilterListItem, error) {
	return s.repo.List(ctx)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*models.Filter, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("%w: invalid id", ErrValidation)
	}
	f, err := s.repo.Get(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, query types.Query) (*models.Filter, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("%w: invalid id", ErrValidation)
	}

	f, err := s.repo.Update(ctx, id, query)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return f, err
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid id", ErrValidation)
	}
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

func (s *service) Apply(ctx context.Context, id uuid.UUID) (types.Query, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("%w: invalid id", ErrValidation)
	}

	f, err := s.repo.Get(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	q, err := placeholder.RenderQuery(
		f.Query,
		time.Now().In(s.loc),
		s.loc,
		s.currentUserID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: render template: %s", ErrValidation, err)
	}
	return q, nil
}
