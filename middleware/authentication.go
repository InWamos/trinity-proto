package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/InWamos/trinity-proto/internal/shared/interfaces/auth/client"
)

type AuthenticationMiddleware struct {
	logger     *slog.Logger
	authClient client.AuthClient
}

func NewAuthenticationMiddleware(logger *slog.Logger, authClient client.AuthClient) *AuthenticationMiddleware {
	middlewareLogger := logger.With(slog.String("component", "authentication_middleware"))
	return &AuthenticationMiddleware{logger: middlewareLogger, authClient: authClient}
}

func (middleware *AuthenticationMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from cookie or Authentication header
		token, err := extractToken(r)
		if err != nil {
			middleware.logger.WarnContext(r.Context(), "failed to extract token", slog.String("err", err.Error()))
			respondWithError(w, http.StatusUnauthorized, "missing or invalid token", "missing_token")
			return
		}

		// Validate session token and get user identity
		userIdentity, err := middleware.authClient.ValidateSession(r.Context(), token)
		if err != nil {
			switch err {
			case client.ErrSessionInvalid:
				middleware.logger.WarnContext(r.Context(), "invalid session", slog.String("token", token))
				respondWithError(w, http.StatusUnauthorized, "invalid session", "invalid_token")

			case client.ErrSessionExpired:
				middleware.logger.WarnContext(r.Context(), "session expired", slog.String("token", token))
				respondWithError(w, http.StatusUnauthorized, "session expired", "expired_token")

			case client.ErrSessionRevoked:
				middleware.logger.WarnContext(r.Context(), "session revoked", slog.String("token", token))
				respondWithError(w, http.StatusUnauthorized, "session revoked", "revoked_token")

			default:
				middleware.logger.ErrorContext(
					r.Context(),
					"unexpected error during session validation",
					slog.Any("err", err),
				)
				respondWithError(w, http.StatusInternalServerError, "authentication failed", "")
				return
			}
			return
		}

		middleware.logger.DebugContext(r.Context(), "session validated successfully",
			slog.String("method", r.Method),
			slog.String("user_id", userIdentity.UserID.String()),
			slog.String("user_role", string(userIdentity.UserRole)),
			slog.String("uri", r.RequestURI))

		// add idp to the context
		ctx := context.WithValue(r.Context(), "IdentityProvider", &userIdentity)

		// Call the next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractToken extracts the session token from the Authorization header
// Expected format: Authorization: Bearer {token}
func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrMissingToken
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrMissingToken
	}

	if parts[1] == "" {
		return "", ErrMissingToken
	}

	return parts[1], nil
}

// JSON error response with WWW-Authenticate header for 401
func respondWithError(w http.ResponseWriter, statusCode int, message string, errorCode string) {
	w.Header().Set("Content-Type", "application/json")

	// RFC 7235
	if statusCode == http.StatusUnauthorized {
		if errorCode != "" {
			w.Header().
				Set("WWW-Authenticate", `Bearer realm="application", error="`+errorCode+`", error_description="`+message+`"`)
		} else {
			w.Header().Set("WWW-Authenticate", `Bearer realm="application"`)
		}
	}

	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
