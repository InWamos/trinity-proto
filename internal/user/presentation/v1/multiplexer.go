package v1

import (
	"net/http"

	"github.com/InWamos/trinity-proto/internal/user/presentation/v1/handlers"
)

type UserMuxV1 struct {
	mux *http.ServeMux
}

func NewUserMuxV1(
	createUserHandler *handlers.CreateUserHandler,
	getUserHandler *handlers.GetUserHandler,
	promoteUserHandler *handlers.PromoteUserHandler,
	demoteUserHandler *handlers.DemoteUserHandler,
	removeUserHandler *handlers.RemoveUserHandler,
) *UserMuxV1 {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", createUserHandler.ServeHTTP)
	mux.HandleFunc("GET /{id}", getUserHandler.ServeHTTP)
	mux.HandleFunc("DELETE /{id}", removeUserHandler.ServeHTTP)

	// User role management
	mux.HandleFunc("PATCH /{id}/promote", promoteUserHandler.ServeHTTP)
	mux.HandleFunc("PATCH /{id}/demote", demoteUserHandler.ServeHTTP)

	return &UserMuxV1{mux: mux}
}

func (um *UserMuxV1) GetMux() *http.ServeMux {
	return um.mux
}
