package v1

import (
	"os"

	"github.com/diegomagalhaes-dev/go-service/business/web/v1/auth"
	"github.com/diegomagalhaes-dev/go-service/business/web/v1/mid"
	"github.com/diegomagalhaes-dev/go-service/foundation/logger"
	"github.com/diegomagalhaes-dev/go-service/foundation/web"
)

type APIMuxConfig struct {
	Build    string
	Shutdown chan os.Signal
	Log      *logger.Logger
	Auth     *auth.Auth
}

type RouteAdder interface {
	Add(app *web.App, cfg APIMuxConfig)
}

func APIMux(cfg APIMuxConfig, routeAdder RouteAdder) *web.App {
	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())

	routeAdder.Add(app, cfg)

	return app
}
