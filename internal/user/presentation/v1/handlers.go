package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/InWamos/trinity-proto/internal/user/application"
	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/go-playground/validator/v10"
)

type createUserForm struct {
	Username    string `json:"username"     validate:"required,alphanum,min=2,max=32"`
	DisplayName string `json:"display_name" validate:"required,alphanumunicode,min=1,max=64"`
	Password    string `json:"password"     validate:"required,alphanumunicode,min=8,max=64"`
	Role        string `json:"role"         validate:"required,oneof=user admin"`
}

type CreateUserHandler struct {
	interactor *application.CreateUser
	validator  *validator.Validate
	logger     *slog.Logger
}

// NewCreateUserHandler builds a new CreateUserHandler.
func NewCreateUserHandler(
	interactor *application.CreateUser,
	validator *validator.Validate,
	logger *slog.Logger,
) *CreateUserHandler {
	cuhLogger := logger.With(slog.String("component", "handler"), slog.String("name", "create_user"))
	return &CreateUserHandler{interactor: interactor, validator: validator, logger: cuhLogger}
}

// ServeHTTP handles an HTTP request to the /user/create endpoint.
func (handler *CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var userForm createUserForm
	if err := json.NewDecoder(r.Body).Decode(&userForm); err != nil {
		handler.logger.DebugContext(r.Context(), "failed to validate json body of request", slog.Any("err", err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	if err := handler.validator.Struct(userForm); err != nil {
		handler.logger.DebugContext(r.Context(), "failed to validate body of request", slog.Any("err", err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	requestDTO := application.CreateUserRequest{
		Username:    userForm.Username,
		DisplayName: userForm.DisplayName,
		Password:    userForm.Password,
		Role:        domain.Role(userForm.Role),
	}
	err := handler.interactor.Execute(r.Context(), requestDTO)
	if err != nil {
		handler.logger.DebugContext(r.Context(), "failed to call the interactor", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).
			Encode(map[string]string{"error": "The server was unable to complete your request. Please try again later"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "The user has been created. you can login now"})
}
