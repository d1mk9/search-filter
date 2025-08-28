package http

import (
	"search-filter/pkg/handlers"
	"search-filter/pkg/service"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API, svc service.Filters) {
	h := handlers.NewFiltersHandler(svc)

	huma.Post(api, "/filters", h.Create, func(op *huma.Operation) {
		op.Description = "Create a saved search filter."
	})

	huma.Get(api, "/filters", h.List, func(op *huma.Operation) {
		op.Description = "List saved filters (no timestamps in items)."
	})

	huma.Get(api, "/filters/{id}", h.Get, func(op *huma.Operation) {
		op.Description = "Get a filter by ID (includes created_at)."
	})

	huma.Put(api, "/filters/{id}", h.Update, func(op *huma.Operation) {
		op.Description = "Update an existing filter."
	})

	huma.Delete(api, "/filters/{id}", h.Delete, func(op *huma.Operation) {
		op.Description = "Delete a filter by ID (204 No Content)."
	})

	huma.Get(api, "/filters/{id}/apply", h.Apply, func(op *huma.Operation) {
		op.Description = "Resolve placeholders and return a ready-to-use query."
	})
}
