package v1

import (
	"net/http"

	"github.com/InWamos/trinity-proto/internal/auth/presentation/v1/handlers"
	"github.com/InWamos/trinity-proto/middleware"
)

type AuthMuxV1 struct {
	mux *http.ServeMux
}

func NewAuthMuxV1(
	loginHandler *handlers.LoginHandler,
	authMiddleware *middleware.AuthenticationMiddleware,
) *AuthMuxV1 {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", loginHandler.ServeHTTP)
	return &AuthMuxV1{mux: mux}
}

func (am *AuthMuxV1) GetMux() *http.ServeMux {
	return am.mux
}
