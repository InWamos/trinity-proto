package v1

import (
	"net/http"

	"github.com/InWamos/trinity-proto/internal/user/presentation/v1/handlers"
)

type UserMuxV1 struct {
	mux *http.ServeMux
}

func NewUserMuxV1(createUserHandler *handlers.CreateUserHandler) *UserMuxV1 {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", createUserHandler.ServeHTTP)

	return &UserMuxV1{mux: mux}
}

func (um *UserMuxV1) GetMux() *http.ServeMux {
	return um.mux
}
