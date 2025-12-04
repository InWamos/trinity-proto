package middleware

import (
	"net/http"
	"time"

	"github.com/InWamos/trinity-proto/config"
	"github.com/rs/cors"
)

type GlobalCORSMiddleware struct {
	allowedOrigin string
}

func NewGlobalCORSMiddleware(serverConfig *config.ServerConfig) *GlobalCORSMiddleware {
	return &GlobalCORSMiddleware{allowedOrigin: serverConfig.AllowedOrigin}
}

func (middleware *GlobalCORSMiddleware) Handler() func(http.Handler) http.Handler {
	corsHeaders := cors.New(cors.Options{
		AllowedOrigins:   []string{middleware.allowedOrigin},
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
