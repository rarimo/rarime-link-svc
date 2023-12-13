package api

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"time"

	"github.com/go-chi/chi"
	"github.com/rarimo/dashboard-rarime-link-svc/internal/config"
	"github.com/rarimo/dashboard-rarime-link-svc/internal/services/api/handlers"
	"gitlab.com/distributed_lab/ape"
)

func Run(ctx context.Context, cfg config.Config) {
	r := chi.NewRouter()

	const slowRequestDurationThreshold = time.Second
	ape.DefaultMiddlewares(r, cfg.Log(), slowRequestDurationThreshold)

	r.Use(
		ape.CtxMiddleware(
			handlers.CtxLog(cfg.Log()),
		),
	)

	r.Route("/v1", func(r chi.Router) {
	})

	cfg.Log().WithFields(logan.F{
		"service": "api",
		"addr":    cfg.Listener().Addr(),
	}).Info("starting api")

	ape.Serve(ctx, r, cfg, ape.ServeOpts{})
}
