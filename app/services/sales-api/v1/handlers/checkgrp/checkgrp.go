// Package checkgrp provides handlers for health check endpoints such as readiness
// and liveness, ensuring that the application is operating correctly and responding to requests.
package checkgrp

import (
	"context"
	"net/http"
	"os"
	"time"

	db "github.com/diegomagalhaes-dev/go-service/business/data/dbsql/pgx"
	"github.com/diegomagalhaes-dev/go-service/foundation/logger"
	"github.com/diegomagalhaes-dev/go-service/foundation/web"
	"github.com/jmoiron/sqlx"
)

// Handlers manages the set of health check endpoints.
type Handlers struct {
	log   *logger.Logger
	build string
	db    *sqlx.DB
}

// New constructs a new Handlers instance for the health check endpoints.
func New(build string, log *logger.Logger, db *sqlx.DB) *Handlers {
	return &Handlers{
		build: build,
		log:   log,
		db:    db,
	}
}

// Readiness checks if the service is ready to accept requests by verifying the database connection.
func (h *Handlers) Readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	status := "ok"
	statusCode := http.StatusOK

	if err := db.StatusCheck(ctx, h.db); err != nil {
		status = "db not ready"
		statusCode = http.StatusInternalServerError
		h.log.Info(ctx, "readiness failure", "status", status)
	}

	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}
	h.log.Info(ctx, "readiness", "status", status)
	return web.Respond(ctx, w, data, statusCode)
}

// Liveness checks if the service is alive and able to respond to requests.
func (h *Handlers) Liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := struct {
		Status     string `json:"status,omitempty"`
		Build      string `json:"build,omitempty"`
		Host       string `json:"host,omitempty"`
		Name       string `json:"name,omitempty"`
		PodIP      string `json:"podIP,omitempty"`
		Node       string `json:"node,omitempty"`
		Namespace  string `json:"namespace,omitempty"`
		GOMAXPROCS string `json:"GOMAXPROCS,omitempty"`
	}{
		Status:     "up",
		Build:      h.build,
		Host:       host,
		Name:       os.Getenv("KUBERNETES_NAME"),
		PodIP:      os.Getenv("KUBERNETES_POD_IP"),
		Node:       os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace:  os.Getenv("KUBERNETES_NAMESPACE"),
		GOMAXPROCS: os.Getenv("GOMAXPROCS"),
	}
	h.log.Info(ctx, "liveness", "status", "OK")

	return web.Respond(ctx, w, data, http.StatusOK)
}
