package handlers

import (
	"net/http"

	"github.com/rarimo/rarime-auth-svc/pkg/auth"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func AuthMiddleware(auth *auth.Client, log *logan.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := auth.ValidateJWT(r.Header)
			if err != nil {
				log.WithError(err).Error("failed to execute auth validate request")
				ape.Render(w, problems.InternalError())
				return
			}

			next.ServeHTTP(w, r.WithContext(CtxUserClaim(claims)(r.Context())))
		})
	}
}

func OptAuthMiddleware(auth *auth.Client, log *logan.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := auth.ValidateJWT(r.Header)
			if err != nil {
				log.WithError(err).Debug("Auth failed, public logic will be executed")
			}

			next.ServeHTTP(w, r.WithContext(CtxUserClaim(claims)(r.Context())))
		})
	}
}
