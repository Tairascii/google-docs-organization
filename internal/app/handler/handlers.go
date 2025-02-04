package handler

import (
	"encoding/json"
	"errors"
	app "github.com/Tairascii/google-docs-organization/internal"
	"github.com/Tairascii/google-docs-organization/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

var (
	ErrUnexpected     = errors.New("unexpected error")
	ErrInvalidRequest = errors.New("invalid request")
)

type Handler struct {
	DI *app.DI
}

func NewHandler(DI *app.DI) *Handler {
	return &Handler{DI: DI}
}

func (h *Handler) InitHandlers() *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))
	r.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			v1.Mount("/organization", handlers(h))
		})
	})
	return r
}

func handlers(h *Handler) http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Post("/", h.Create)
	})

	return rg
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var payload CreateOrgPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		pkg.JSONErrorResponseWriter(w, ErrInvalidRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	id, err := h.DI.UseCase.Org.CreateOrg(ctx, payload.Title)
	if err != nil {
		pkg.JSONErrorResponseWriter(w, ErrUnexpected, http.StatusInternalServerError)
		return
	}

	pkg.JSONResponseWriter[CreateOrgResponse](w, CreateOrgResponse{ID: id}, http.StatusOK)
}
