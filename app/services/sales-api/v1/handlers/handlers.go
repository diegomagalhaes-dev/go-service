package handlers

import (
	"github.com/diegomagalhaes-dev/go-service/app/services/sales-api/v1/handlers/hackgrp"
	v1 "github.com/diegomagalhaes-dev/go-service/business/web/v1"
	"github.com/diegomagalhaes-dev/go-service/foundation/web"
)

type Routes struct{}

// Add implements the RouterAdder interface.
func (Routes) Add(app *web.App, cfg v1.APIMuxConfig) {
	hackgrp.Routes(app)
}
