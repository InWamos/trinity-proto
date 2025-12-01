package main

import (
	"net/http"

	"github.com/InWamos/trinity-proto/setup"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(setup.NewHTTPServer),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
