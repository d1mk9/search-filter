// pkg/handlers/handlers.go
package handlers

import (
	"context"
	"errors"
	"time"

	"search-filter/pkg/models"
	"search-filter/pkg/service"
	"search-filter/pkg/types"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type FiltersHandler struct {
	svc service.Filters
}

func NewFiltersHandler(svc service.Filters) *FiltersHandler {
	return &FiltersHandler{svc: svc}
}

type FilterDTO struct {
	ID        uuid.UUID   `json:"id"`
	Name      string      `json:"name"`
	Query     types.Query `json:"query"`
	CreatedAt time.Time   `json:"created_at"`
}

type FilterListItemDTO struct {
	ID    uuid.UUID   `json:"id"`
	Name  string      `json:"name"`
	Query types.Query `json:"query"`
}

func toFilterDTO(m models.Filter) FilterDTO {
	return FilterDTO{
		ID:        m.ID,
		Name:      m.Name,
		Query:     m.Query,
		CreatedAt: m.CreatedAt,
	}
}

type createFilterBody struct {
	Name  string      `json:"name" minLength:"1"`
	Query types.Query `json:"query" jsonschema:"minProperties=1"`
}
type createFilterInput struct {
	Body createFilterBody `json:"body"`
}
type createFilterOutput struct {
	Body FilterDTO `json:"body"`
}

func (h *FiltersHandler) Create(ctx context.Context, in *createFilterInput) (*createFilterOutput, error) {
	f, err := h.svc.Create(ctx, in.Body.Name, in.Body.Query)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrValidation):
			return nil, huma.Error422UnprocessableEntity(err.Error())
		default:
			return nil, huma.Error500InternalServerError("internal error")
		}
	}
	return &createFilterOutput{Body: toFilterDTO(*f)}, nil
}

type listFiltersOutput struct {
	Body []FilterListItemDTO `json:"body"`
}

func (h *FiltersHandler) List(ctx context.Context, _ *struct{}) (*listFiltersOutput, error) {
	items, err := h.svc.List(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("internal error")
	}
	out := make([]FilterListItemDTO, 0, len(items))
	for _, it := range items {
		out = append(out, FilterListItemDTO{
			ID:    it.ID,
			Name:  it.Name,
			Query: it.Query,
		})
	}
	return &listFiltersOutput{Body: out}, nil
}

type IdPath struct {
	ID uuid.UUID `path:"id" format:"uuid"`
}
type getFilterInput struct {
	IdPath
}
type getFilterOutput struct {
	Body FilterDTO `json:"body"`
}

func (h *FiltersHandler) Get(ctx context.Context, in *getFilterInput) (*getFilterOutput, error) {
	f, err := h.svc.Get(ctx, in.ID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			return nil, huma.Error404NotFound("not found")
		case errors.Is(err, service.ErrValidation):
			return nil, huma.Error422UnprocessableEntity(err.Error())
		default:
			return nil, huma.Error500InternalServerError("internal error")
		}
	}
	return &getFilterOutput{Body: toFilterDTO(*f)}, nil
}

type updateFilterBody struct {
	Query types.Query `json:"query" jsonschema:"minProperties=1"`
}
type updateFilterInput struct {
	IdPath
	Body updateFilterBody `json:"body"`
}
type updateFilterOutput struct {
	Body FilterDTO `json:"body"`
}

func (h *FiltersHandler) Update(ctx context.Context, in *updateFilterInput) (*updateFilterOutput, error) {
	f, err := h.svc.Update(ctx, in.ID, in.Body.Query)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			return nil, huma.Error404NotFound("not found")
		case errors.Is(err, service.ErrValidation):
			return nil, huma.Error422UnprocessableEntity(err.Error())
		default:
			return nil, huma.Error500InternalServerError("internal error")
		}
	}
	return &updateFilterOutput{Body: toFilterDTO(*f)}, nil
}

type deleteFilterInput struct {
	IdPath
}

func (h *FiltersHandler) Delete(ctx context.Context, in *deleteFilterInput) (*struct{}, error) {
	if err := h.svc.Delete(ctx, in.ID); err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			return nil, huma.Error404NotFound("not found")
		case errors.Is(err, service.ErrValidation):
			return nil, huma.Error422UnprocessableEntity(err.Error())
		default:
			return nil, huma.Error500InternalServerError("internal error")
		}
	}
	return nil, nil
}

type applyFilterInput struct {
	IdPath
}
type applyFilterOutput struct {
	Body types.Query `json:"body"`
}

func (h *FiltersHandler) Apply(ctx context.Context, in *applyFilterInput) (*applyFilterOutput, error) {
	q, err := h.svc.Apply(ctx, in.ID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			return nil, huma.Error404NotFound("not found")
		case errors.Is(err, service.ErrValidation):
			return nil, huma.Error422UnprocessableEntity(err.Error())
		default:
			return nil, huma.Error500InternalServerError("internal error")
		}
	}
	return &applyFilterOutput{Body: q}, nil
}
