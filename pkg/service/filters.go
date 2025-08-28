package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"search-filter/pkg/models"
	"search-filter/pkg/placeholder"
	"search-filter/pkg/repository"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrValidation = errors.New("validation error")
	ErrConflict   = errors.New("conflict")
)

type Filters interface {
	Create(ctx context.Context, name string, query Query) (*models.Filter, error)
	List(ctx context.Context) ([]models.FilterListItem, error)
	Get(ctx context.Context, id int64) (*models.Filter, error)
	Update(ctx context.Context, id int64, query Query) (*models.Filter, error)
	Delete(ctx context.Context, id int64) error
	Apply(ctx context.Context, id int64) (Query, error)
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

func (s *service) Create(ctx context.Context, name string, query Query) (*models.Filter, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrValidation)
	}
	if err := validateQueryObjectMap(query); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidation, err)
	}
	raw, err := marshalQuery(query)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidation, err)
	}
	return s.repo.Create(ctx, name, raw)
}

func (s *service) List(ctx context.Context) ([]models.FilterListItem, error) {
	return s.repo.List(ctx)
}

func (s *service) Get(ctx context.Context, id int64) (*models.Filter, error) {
	if id <= 0 {
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

func (s *service) Update(ctx context.Context, id int64, query Query) (*models.Filter, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid id", ErrValidation)
	}
	if err := validateQueryObjectMap(query); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidation, err)
	}
	raw, err := marshalQuery(query)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidation, err)
	}
	f, err := s.repo.Update(ctx, id, raw)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: invalid id", ErrValidation)
	}
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Apply(ctx context.Context, id int64) (Query, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid id", ErrValidation)
	}
	f, err := s.repo.Get(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	now := time.Now().In(s.loc)
	resolvedRaw, err := placeholder.RenderTemplate(
		f.Query,
		placeholder.TemplateCtx{
			Now:         now,
			Loc:         s.loc,
			CurrentUser: s.currentUserID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w: render template: %s", ErrValidation, err)
	}

	q, err := unmarshalQuery(resolvedRaw)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidation, err)
	}
	return q, nil
}
