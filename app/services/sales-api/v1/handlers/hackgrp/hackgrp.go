package hackgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"github.com/diegomagalhaes-dev/go-service/business/web/v1/response"
	"github.com/diegomagalhaes-dev/go-service/foundation/web"
)

func Hack(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100) % 2; n == 0 {
		panic("OH MY GOODNESS WE PANIC'D")
		return response.NewError(errors.New("TRUST ERROR"), http.StatusBadRequest)
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
