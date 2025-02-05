package handler

import (
	"encoding/json"
	"errors"
	app "github.com/Tairascii/google-docs-organization/internal"
	"github.com/Tairascii/google-docs-organization/internal/app/usecase"
	"github.com/Tairascii/google-docs-organization/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"net/http"
)

var (
	ErrUnexpected     = errors.New("unexpected error")
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidUserId  = errors.New("invalid user id")
	ErrInvalidOrgId   = errors.New("invalid org id")
	ErrInvalidOwnerId = errors.New("invalid owner id")
	ErrAuth           = errors.New("authentication failed")
)

// TODO move to apigw and use vault
const (
	accessSecret = "yoS0baK1Ya"
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
	r.Use(ParseToken(accessSecret))
	r.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			v1.Mount("/organization", handlers(h))
			v1.Mount("/users", usersHandlers(h))
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

	if payload.Title == "" {
		pkg.JSONErrorResponseWriter(w, ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	id, err := h.DI.UseCase.Org.CreateOrg(ctx, payload.Title)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidOwnerId) {
			pkg.JSONErrorResponseWriter(w, ErrInvalidOwnerId, http.StatusBadRequest)
			return
		}
		pkg.JSONErrorResponseWriter(w, ErrUnexpected, http.StatusInternalServerError)
		return
	}

	pkg.JSONResponseWriter[CreateOrgResponse](w, CreateOrgResponse{ID: id.String()}, http.StatusOK)
}

func usersHandlers(h *Handler) http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Post("/add", h.AddUser)
	})

	return rg
}

func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	var payload AddUserPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		pkg.JSONErrorResponseWriter(w, ErrInvalidRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	userId, err := uuid.Parse(payload.UserId)
	if err != nil {
		pkg.JSONErrorResponseWriter(w, ErrInvalidUserId, http.StatusBadRequest)
		return
	}

	orgId, err := uuid.Parse(payload.OrgId)
	if err != nil {
		pkg.JSONErrorResponseWriter(w, ErrInvalidOrgId, http.StatusBadRequest)
		return
	}

	if payload.Role == "" {
		pkg.JSONErrorResponseWriter(w, ErrInvalidOrgId, http.StatusBadRequest)
		return
	}

	err = h.DI.UseCase.Org.AddUser(ctx, orgId, userId, payload.Role)
	if err != nil {
		pkg.JSONErrorResponseWriter(w, ErrUnexpected, http.StatusInternalServerError)
		return
	}

	pkg.EmptyResponseWriter(w, http.StatusNoContent)
}
