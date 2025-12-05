package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/InWamos/trinity-proto/config"
)

type TrustedProxyMiddleware struct {
	trustedProxies []string
}

func NewTrustedProxyMiddleware(serverConfig *config.ServerConfig) *TrustedProxyMiddleware {
	return &TrustedProxyMiddleware{trustedProxies: []string{serverConfig.TrustedProxy}}
}

func (middleware *TrustedProxyMiddleware) Handler(next http.Handler) http.Handler {
	trustedMap := make(map[string]bool)
	for _, ip := range middleware.trustedProxies {
		trustedMap[ip] = true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the remote IP (without port)
		remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
		// Only trust X-Forwarded-For if request comes from trusted proxy
		if trustedMap[remoteIP] {
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				// Take the first IP (leftmost = original client)
				if idx := strings.Index(xff, ","); idx != -1 {
					xff = strings.TrimSpace(xff[:idx])
				}
				// Store in context or custom header for handlers to use
				r.RemoteAddr = xff + ":0"
			}
		}
		// If not from trusted proxy, keep original RemoteAddr
		next.ServeHTTP(w, r)
	})
}
