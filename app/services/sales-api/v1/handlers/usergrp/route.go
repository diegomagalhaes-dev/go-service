package usergrp

import (
	"net/http"

	"github.com/diegomagalhaes-dev/go-service/business/core/event"
	"github.com/diegomagalhaes-dev/go-service/business/core/user"
	"github.com/diegomagalhaes-dev/go-service/business/core/user/stores/userdb"
	db "github.com/diegomagalhaes-dev/go-service/business/data/dbsql/pgx"
	"github.com/diegomagalhaes-dev/go-service/business/web/v1/auth"
	"github.com/diegomagalhaes-dev/go-service/business/web/v1/mid"
	"github.com/diegomagalhaes-dev/go-service/foundation/logger"
	"github.com/diegomagalhaes-dev/go-service/foundation/web"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
	Auth  *auth.Auth
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.Auth)
	ruleAdmin := mid.Authorize(cfg.Auth, auth.RuleAdminOnly)
	ruleAdminOrSubject := mid.Authorize(cfg.Auth, auth.RuleAdminOrSubject)
	tran := mid.ExecuteInTransation(cfg.Log, db.NewBeginner(cfg.DB))

	envCore := event.NewCore(cfg.Log)
	usrCore := user.NewCore(cfg.Log, envCore, userdb.NewStore(cfg.Log, cfg.DB))

	hdl := New(usrCore, cfg.Auth)
	app.Handle(http.MethodGet, version, "/users/token/:kid", hdl.Token)
	app.Handle(http.MethodGet, version, "/users", hdl.Query, authen, ruleAdmin)
	app.Handle(http.MethodGet, version, "/users/:user_id", hdl.QueryByID, authen, ruleAdminOrSubject)
	app.Handle(http.MethodPost, version, "/users", hdl.Create, tran)
	app.Handle(http.MethodPut, version, "/users/:user_id", hdl.Update, authen, ruleAdminOrSubject, tran)
	app.Handle(http.MethodDelete, version, "/users/:user_id", hdl.Delete, authen, ruleAdminOrSubject, tran)
}
