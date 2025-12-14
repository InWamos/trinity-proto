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

	// User CRUD operations with full paths
	mux.HandleFunc("POST /api/v1/users", createUserHandler.ServeHTTP)
	mux.HandleFunc("GET /api/v1/users/{id}", getUserHandler.ServeHTTP)
	mux.HandleFunc("DELETE /api/v1/users/{id}", removeUserHandler.ServeHTTP)
	
	// User role management
	mux.HandleFunc("PATCH /api/v1/users/{id}/promote", promoteUserHandler.ServeHTTP)
	mux.HandleFunc("PATCH /api/v1/users/{id}/demote", demoteUserHandler.ServeHTTP)

	return &UserMuxV1{mux: mux}
}

func (um *UserMuxV1) GetMux() *http.ServeMux {
	return um.mux
}
