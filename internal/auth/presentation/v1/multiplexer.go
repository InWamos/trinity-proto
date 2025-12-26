package v1

import (
	"github.com/InWamos/trinity-proto/internal/auth/presentation/v1/handlers"
	"github.com/go-chi/chi/v5"
)

type AuthMuxV1 struct {
	mux *chi.Mux
}

func NewAuthMuxV1(
	loginHandler *handlers.LoginHandler,
) *AuthMuxV1 {
	mux := chi.NewRouter()
	mux.Post("/login", loginHandler.ServeHTTP)
	return &AuthMuxV1{mux: mux}
}

func (am *AuthMuxV1) GetMux() *chi.Mux {
	return am.mux
}
