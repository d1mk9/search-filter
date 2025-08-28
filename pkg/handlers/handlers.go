package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"search-filter/pkg/models"
	"search-filter/pkg/service"

	"github.com/danielgtaylor/huma/v2"
)

type FiltersHandler struct {
	svc service.Filters
}

func NewFiltersHandler(svc service.Filters) *FiltersHandler {
	return &FiltersHandler{svc: svc}
}

type FilterDTO struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	Query     json.RawMessage `json:"query"`
	CreatedAt string          `json:"created_at"`
}

type FilterListItemDTO struct {
	ID    int64           `json:"id"`
	Name  string          `json:"name"`
	Query json.RawMessage `json:"query"`
}

func toFilterDTO(m models.Filter) FilterDTO {
	return FilterDTO{
		ID:        m.ID,
		Name:      m.Name,
		Query:     m.Query,
		CreatedAt: m.CreatedAt.UTC().Format(time.RFC3339),
	}
}

type createFilterBody struct {
	Name  string          `json:"name"`
	Query json.RawMessage `json:"query"`
}
type createFilterInput struct {
	Body createFilterBody `json:"body"`
}
type createFilterOutput struct {
	Body FilterDTO `json:"body"`
}

func (h *FiltersHandler) Create(ctx context.Context, in *createFilterInput) (*createFilterOutput, error) {
	if in == nil || in.Body.Name == "" || len(in.Body.Query) == 0 {
		return nil, huma.Error400BadRequest("name and query are required")
	}

	var q service.Query
	if err := json.Unmarshal(in.Body.Query, &q); err != nil {
		return nil, huma.Error400BadRequest("query must be a valid JSON object")
	}

	f, err := h.svc.Create(ctx, in.Body.Name, q)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrValidation):
			return nil, huma.Error400BadRequest(err.Error())
		case errors.Is(err, service.ErrConflict):
			return nil, huma.Error409Conflict("conflict")
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
	ID int64 `path:"id"`
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
			return nil, huma.Error400BadRequest(err.Error())
		default:
			return nil, huma.Error500InternalServerError("internal error")
		}
	}
	return &getFilterOutput{Body: toFilterDTO(*f)}, nil
}

type updateFilterBody struct {
	Query json.RawMessage `json:"query"`
}
type updateFilterInput struct {
	IdPath
	Body updateFilterBody `json:"body"`
}
type updateFilterOutput struct {
	Body FilterDTO `json:"body"`
}

func (h *FiltersHandler) Update(ctx context.Context, in *updateFilterInput) (*updateFilterOutput, error) {
	if len(in.Body.Query) == 0 {
		return nil, huma.Error400BadRequest("query is required")
	}

	var q service.Query
	if err := json.Unmarshal(in.Body.Query, &q); err != nil {
		return nil, huma.Error400BadRequest("query must be a valid JSON object")
	}

	f, err := h.svc.Update(ctx, in.ID, q)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			return nil, huma.Error404NotFound("not found")
		case errors.Is(err, service.ErrValidation):
			return nil, huma.Error400BadRequest(err.Error())
		case errors.Is(err, service.ErrConflict):
			return nil, huma.Error409Conflict("conflict")
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
			return nil, huma.Error400BadRequest(err.Error())
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
	Body json.RawMessage `json:"body"`
}

func (h *FiltersHandler) Apply(ctx context.Context, in *applyFilterInput) (*applyFilterOutput, error) {
	q, err := h.svc.Apply(ctx, in.ID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			return nil, huma.Error404NotFound("not found")
		case errors.Is(err, service.ErrValidation):
			return nil, huma.Error400BadRequest(err.Error())
		default:
			return nil, huma.Error500InternalServerError("internal error")
		}
	}

	raw, err := json.Marshal(q)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to encode response")
	}
	return &applyFilterOutput{Body: raw}, nil
}
