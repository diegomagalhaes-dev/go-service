package checkgrp

import (
	"net/http"

	"github.com/diegomagalhaes-dev/go-service/foundation/logger"
	"github.com/diegomagalhaes-dev/go-service/foundation/web"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	hdl := New(cfg.Build, cfg.Log, cfg.DB)
	app.HandleNoMiddleware(http.MethodGet, version, "/readiness", hdl.Readiness)
	app.HandleNoMiddleware(http.MethodGet, version, "/liveness", hdl.Liveness)
}
