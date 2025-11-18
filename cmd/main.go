package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/InWamos/trinity-proto/config"
	"github.com/gin-contrib/cors"
	ginslog "github.com/gin-contrib/slog"
	"github.com/gin-gonic/gin"
)

func getConfig() (*config.AppConfig, error) {
	config, err := config.NewAppConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func bootstrapServer(server *gin.Engine, allowOrigin string, trustedProxy string) {
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowOrigin},
		AllowMethods:     []string{"GET", "DELETE", "PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	server.Use(ginslog.SetLogger(
		// Change log writer
		ginslog.WithWriter(os.Stdout),
		// Use UTC timestamps
		ginslog.WithUTC(true),
		// Skip health check and static routes
		ginslog.WithSkipPath([]string{"/healthz", "/metrics"}),
		// Change default log levels
		ginslog.WithDefaultLevel(slog.LevelDebug),
		ginslog.WithClientErrorLevel(slog.LevelWarn),
		ginslog.WithServerErrorLevel(slog.LevelError),
		// Log message customization
		ginslog.WithMessage("Handled request"),
		// Set specific log level for a given path
		ginslog.WithPathLevel(map[string]slog.Level{"/foo": slog.LevelInfo}),
		// Inject user agent
		ginslog.WithContext(func(c *gin.Context, rec *slog.Record) *slog.Record {
			rec.Add("user_agent", c.Request.UserAgent())
			return rec
		}),
		// Provide your own logger (to add global fields, etc.)
		ginslog.WithLogger(func(c *gin.Context, l *slog.Logger) *slog.Logger {
			return l.With("request_id", c.GetString("request_id"))
		}),
		// Custom Skipper (function: skip logging if ...), example:
		ginslog.WithSkipper(func(c *gin.Context) bool {
			return c.Request.Method == http.MethodOptions
		}),
		// Hide sensitive request headers from logs (optional, default hides Authorization/Cookie/etc)
		ginslog.WithRequestHeader(true), // Set to false to disable logging all request headers
		ginslog.WithHiddenRequestHeaders([]string{
			"authorization", "cookie", "x-csrf-token", "bearer", // set your own or reset
		}),
	))
	server.Use(gin.Recovery())
	server.SetTrustedProxies([]string{trustedProxy})
	server.RemoteIPHeaders = []string{"X-Forwarded-For"}
}

func runServer(server *gin.Engine, bindIp string, port int) {
	bindAddress := bindIp + ":" + string(rune(port))
	if err := server.Run(bindAddress); err != nil {
		slog.Error("Failed to start server", "error", err, "bind_address", bindAddress)
		panic(err)
	}
}

func main() {
	config, err := getConfig()
	if err != nil {
		slog.Error("Failed to parse config", "error", err)
		os.Exit(1)
	}
	server := gin.New()
	bootstrapServer(server, config.GinConfig.AllowedOrigin, config.GinConfig.TrustedProxy)
	runServer(server, config.GinConfig.BindAddress, config.GinConfig.Port)
}
