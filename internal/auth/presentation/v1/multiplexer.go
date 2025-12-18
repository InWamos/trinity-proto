package v1

import (
	"net/http"

	"github.com/InWamos/trinity-proto/internal/auth/presentation/v1/handlers"
)

type AuthMuxV1 struct {
	mux *http.ServeMux
}

func NewUserMuxV1(
	loginHandler *handlers.LoginHandler,
) *AuthMuxV1 {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", loginHandler.ServeHTTP)
	return &AuthMuxV1{mux: mux}
}

func (um *AuthMuxV1) GetMux() *http.ServeMux {
	return um.mux
}
