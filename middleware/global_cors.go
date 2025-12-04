package middleware

import (
	"net/http"
	"time"

	"github.com/rs/cors"
)

func GlobalCORSHeaders(allowedOrigin []string) func(http.Handler) http.Handler {
	corsHeaders := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigin,
		AllowedMethods:   []string{"GET", "DELETE", "PUT", "PATCH", "POST"},
		AllowedHeaders:   []string{"Origin"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           int(time.Hour.Seconds() * 24),
	})
	return func(h http.Handler) http.Handler {
		return corsHeaders.Handler(h)
	}
}
