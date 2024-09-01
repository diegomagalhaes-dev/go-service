package handlers

import (
	"github.com/diegomagalhaes-dev/go-service/app/services/sales-api/v1/handlers/checkgrp"
	"github.com/diegomagalhaes-dev/go-service/app/services/sales-api/v1/handlers/hackgrp"
	"github.com/diegomagalhaes-dev/go-service/app/services/sales-api/v1/handlers/productgrp"
	"github.com/diegomagalhaes-dev/go-service/app/services/sales-api/v1/handlers/usergrp"
	"github.com/diegomagalhaes-dev/go-service/app/services/sales-api/v1/handlers/usersummarygrp"
	v1 "github.com/diegomagalhaes-dev/go-service/business/web/v1"
	"github.com/diegomagalhaes-dev/go-service/foundation/web"
)

type Routes struct{}

// Add implements the RouterAdder interface to add all routes.
func (Routes) Add(app *web.App, apiCfg v1.APIMuxConfig) {
	hackgrp.Routes(app, hackgrp.Config{
		Auth: apiCfg.Auth,
	})

	checkgrp.Routes(app, checkgrp.Config{
		Build: apiCfg.Build,
		Log:   apiCfg.Log,
		DB:    apiCfg.DB,
	})

	usergrp.Routes(app, usergrp.Config{
		Build: apiCfg.Build,
		Log:   apiCfg.Log,
		DB:    apiCfg.DB,
		Auth:  apiCfg.Auth,
	})

	productgrp.Routes(app, productgrp.Config{
		Build: apiCfg.Build,
		Log:   apiCfg.Log,
		DB:    apiCfg.DB,
		Auth:  apiCfg.Auth,
	})

	usersummarygrp.Routes(app, usersummarygrp.Config{
		Log:  apiCfg.Log,
		Auth: apiCfg.Auth,
		DB:   apiCfg.DB,
	})
}
