package productgrp

import (
	"net/http"

	"github.com/diegomagalhaes-dev/go-service/business/core/event"
	"github.com/diegomagalhaes-dev/go-service/business/core/product"
	"github.com/diegomagalhaes-dev/go-service/business/core/product/stores/productdb"
	"github.com/diegomagalhaes-dev/go-service/business/core/user"
	"github.com/diegomagalhaes-dev/go-service/business/core/user/stores/usercache"
	"github.com/diegomagalhaes-dev/go-service/business/core/user/stores/userdb"
	db "github.com/diegomagalhaes-dev/go-service/business/data/dbsql/pgx"
	"github.com/diegomagalhaes-dev/go-service/business/web/v1/auth"
	"github.com/diegomagalhaes-dev/go-service/business/web/v1/mid"
	"github.com/diegomagalhaes-dev/go-service/foundation/logger"
	"github.com/diegomagalhaes-dev/go-service/foundation/web"
	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
	Auth  *auth.Auth
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	envCore := event.NewCore(cfg.Log)
	usrCore := user.NewCore(cfg.Log, envCore, usercache.NewStore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB)))
	prdCore := product.NewCore(cfg.Log, envCore, usrCore, productdb.NewStore(cfg.Log, cfg.DB))

	authen := mid.Authenticate(cfg.Auth)
	tran := mid.ExecuteInTransation(cfg.Log, db.NewBeginner(cfg.DB))

	hdl := New(prdCore, usrCore)
	app.Handle(http.MethodGet, version, "/products", hdl.Query, authen)
	app.Handle(http.MethodGet, version, "/products/:product_id", hdl.QueryByID, authen)
	app.Handle(http.MethodPost, version, "/products", hdl.Create, authen)
	app.Handle(http.MethodPut, version, "/products/:product_id", hdl.Update, authen, tran)
	app.Handle(http.MethodDelete, version, "/products/:product_id", hdl.Delete, authen, tran)
}
