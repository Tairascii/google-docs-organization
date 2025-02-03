package handler

import (
	app "github.com/Tairascii/google-docs-organization/internal"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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
	return r
}
