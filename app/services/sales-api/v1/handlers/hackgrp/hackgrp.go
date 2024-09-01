// Package hackgrp provides a handler for simulating an endpoint that randomly
// returns an error or a success response, used primarily for testing error handling
// and response mechanisms.
package hackgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"github.com/diegomagalhaes-dev/go-service/business/web/v1/response"
	"github.com/diegomagalhaes-dev/go-service/foundation/web"
)

// Hack randomly returns a bad request error or an OK status based on a random number.
// It simulates an endpoint for testing purposes.
func Hack(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100) % 2; n == 0 {
		return response.NewError(errors.New("TRUST ERROR"), http.StatusBadRequest)
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
