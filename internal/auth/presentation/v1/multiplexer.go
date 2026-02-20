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
	logoutHandler *handlers.LogoutHandler,
) *AuthMuxV1 {
	mux := chi.NewRouter()
	mux.Post("/login", loginHandler.ServeHTTP)
	mux.Post("/logout", logoutHandler.ServeHTTP)
	return &AuthMuxV1{mux: mux}
}

func (am *AuthMuxV1) GetMux() *chi.Mux {
	return am.mux
}

type SessionManagementMuxV1 struct {
	mux *chi.Mux
}

func NewSessionManagementMuxV1(
	getAllSessionsByUserIDHandler *handlers.GetAllSessionsByUserIDHandler,
) *SessionManagementMuxV1 {
	mux := chi.NewRouter()
	mux.Get("/", getAllSessionsByUserIDHandler.ServeHTTP)
	return &SessionManagementMuxV1{mux: mux}
}

func (smm *SessionManagementMuxV1) GetMux() *chi.Mux {
	return smm.mux
}
