package v1

import (
	"net/http"

	"github.com/InWamos/trinity-proto/internal/user/presentation/v1/handlers"
)

func GetUserMuxV1(createUserHandler *handlers.CreateUserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", createUserHandler.ServeHTTP)

	return mux
}
