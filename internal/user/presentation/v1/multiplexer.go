package v1

import (
	"github.com/InWamos/trinity-proto/internal/user/presentation/v1/handlers"
	"github.com/go-chi/chi/v5"
)

type UserMuxV1 struct {
	mux *chi.Mux
}

func NewUserMuxV1(
	createUserHandler *handlers.CreateUserHandler,
	getUserHandler *handlers.GetUserHandler,
	promoteUserHandler *handlers.PromoteUserHandler,
	demoteUserHandler *handlers.DemoteUserHandler,
	removeUserHandler *handlers.RemoveUserHandler,
) *UserMuxV1 {
	mux := chi.NewRouter()
	mux.Post("/", createUserHandler.ServeHTTP)
	mux.Get("/{id}", getUserHandler.ServeHTTP)
	mux.Delete("/{id}", removeUserHandler.ServeHTTP)
	mux.Patch("/{id}/promote", promoteUserHandler.ServeHTTP)
	mux.Patch("/{id}/demote", demoteUserHandler.ServeHTTP)
	return &UserMuxV1{mux: mux}
}

func (um *UserMuxV1) GetMux() *chi.Mux {
	return um.mux
}
