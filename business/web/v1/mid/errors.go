package mid

import (
	"context"
	"net/http"

	"github.com/diegomagalhaes-dev/go-service/business/web/v1/auth"
	"github.com/diegomagalhaes-dev/go-service/business/web/v1/response"
	"github.com/diegomagalhaes-dev/go-service/foundation/logger"
	"github.com/diegomagalhaes-dev/go-service/foundation/validate"
	"github.com/diegomagalhaes-dev/go-service/foundation/web"
)

func Errors(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				log.Error(ctx, "message", "msg", err)

				var er response.ErrorDocument
				var status int

				switch {
				case response.IsError(err):
					reqErr := response.GetError(err)

					if validate.IsFieldErrors(reqErr.Err) {
						fieldErrors := validate.GetFieldErrors(reqErr.Err)
						er = response.ErrorDocument{
							Error:  "data validation error",
							Fields: fieldErrors.Fields(),
						}
						status = reqErr.Status
						break
					}

					er = response.ErrorDocument{
						Error: reqErr.Error(),
					}
					status = reqErr.Status

				case auth.IsAuthError(err):
					er = response.ErrorDocument{
						Error: http.StatusText(http.StatusUnauthorized),
					}
					status = http.StatusUnauthorized

				default:
					er = response.ErrorDocument{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
