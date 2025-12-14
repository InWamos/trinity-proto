package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/InWamos/trinity-proto/internal/user/application"
	"github.com/InWamos/trinity-proto/internal/user/domain"
	"github.com/InWamos/trinity-proto/internal/user/presentation/service"
)

// CreateUserResponse represents the response from the CreateUser endpoint
// @Description User creation response with ID
type CreateUserResponse struct {
	Message string `json:"message" example:"The user has been created. You can login now"`
	ID      string `json:"id" example:"019b1a49-dbf6-74d6-97bf-2d7e57d30c75"`
}

// ErrorResponse represents an error response
// @Description Standard error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request body"`
}

type createUserForm struct {
	Username    string `json:"username"     validate:"required,alphanum,min=2,max=32"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=64"`
	Password    string `json:"password"     validate:"required,alphanumunicode,min=8,max=64"`
	Role        string `json:"role"         validate:"required,oneof=user admin"`
}

type CreateUserHandler struct {
	interactor *application.CreateUser
	validator  service.PostFormValidator
	logger     *slog.Logger
}

// NewCreateUserHandler builds a new CreateUserHandler.
func NewCreateUserHandler(
	interactor *application.CreateUser,
	validator service.PostFormValidator,
	logger *slog.Logger,
) *CreateUserHandler {
	cuhLogger := logger.With(slog.String("component", "handler"), slog.String("name", "create_user"))
	return &CreateUserHandler{interactor: interactor, validator: validator, logger: cuhLogger}
}

// ServeHTTP handles an HTTP request to create a user.
// @Summary Create a new user
// @Description Create a new user with username, display name, password and role
// @Tags users
// @Accept json
// @Produce json
// @Param request body createUserForm true "User creation request"
// @Success 201 {object} CreateUserResponse "User created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request (validation failed)"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/users [post]
func (handler *CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userForm createUserForm
	if err := handler.validator.ValidateBody(r.Body, &userForm); err != nil {
		handler.logger.DebugContext(r.Context(), "failed to validate the form", slog.Any("err", err))
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
	response, err := handler.interactor.Execute(r.Context(), requestDTO)
	if err != nil {
		handler.logger.DebugContext(r.Context(), "failed to call the interactor", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).
			Encode(map[string]string{"error": "The server was unable to complete your request. Please try again later"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "The user has been created. You can login now",
		"id":      response.UserID.String(),
	})
}
