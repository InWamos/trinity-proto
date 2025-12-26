package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/InWamos/trinity-proto/internal/auth/application"
	"github.com/InWamos/trinity-proto/internal/user/presentation/service"
)

// LoginResponse represents the response from the Login endpoint
//
//	@Description	Login response with session token
type LoginResponse struct {
	Message string `json:"message" example:"Login successful"`
	Token   string `json:"token"   example:"dGVzdC10b2tlbi0xMjM0NTY3ODkw"`
}

// ErrorResponse represents an error response
//
//	@Description	Standard error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid credentials"`
}

type loginForm struct {
	Username string `json:"username" validate:"required,alphanum,min=2,max=32"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type LoginHandler struct {
	interactor *application.AddSession
	validator  service.PostFormValidator
	logger     *slog.Logger
}

// NewLoginHandler builds a new LoginHandler.
func NewLoginHandler(
	interactor *application.AddSession,
	validator service.PostFormValidator,
	logger *slog.Logger,
) *LoginHandler {
	lhLogger := logger.With(slog.String("component", "handler"), slog.String("name", "login"))
	return &LoginHandler{interactor: interactor, validator: validator, logger: lhLogger}
}

// ServeHTTP handles an HTTP request to login a user.
//
//	@Summary		User login
//	@Description	Authenticate a user with username and password, returns session token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loginForm		true	"Login credentials"
//	@Success		200		{object}	LoginResponse	"Login successful"
//	@Failure		400		{object}	ErrorResponse	"Invalid request (validation failed)"
//	@Failure		401		{object}	ErrorResponse	"Invalid credentials"
//	@Failure		500		{object}	ErrorResponse	"Server error"
//	@Router			/v1/auth/login [post]
func (handler *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var form loginForm
	if err := handler.validator.ValidateBody(r.Body, &form); err != nil {
		handler.logger.DebugContext(r.Context(), "failed to validate the form", slog.Any("err", err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Extract IP address and user agent from request
	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	requestDTO := application.AddSessionRequest{
		Username:  form.Username,
		Password:  form.Password,
		IpAddress: ipAddress,
		UserAgent: userAgent,
	}

	response, err := handler.interactor.Execute(r.Context(), requestDTO)
	if err != nil {
		handler.logger.DebugContext(r.Context(), "failed to execute login", slog.Any("err", err))

		switch {
		case errors.Is(err, application.ErrInvalidCredentials):
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "The server was unable to complete your request. Please try again later",
			})
		}
		return
	}

	// Set session token as HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    response.Session.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(response.Session.ExpiresAt.Sub(response.Session.CreatedAt).Seconds()),
	})

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   response.Session.Token,
	})
}
