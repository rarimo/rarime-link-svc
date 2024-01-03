package api

import (
	"context"
	"github.com/rarimo/rarime-link-svc/internal/services/proofs_cleaner"
	"gitlab.com/distributed_lab/logan/v3"
	"time"

	"github.com/go-chi/chi"
	"github.com/rarimo/rarime-link-svc/internal/config"
	"github.com/rarimo/rarime-link-svc/internal/services/api/handlers"
	"gitlab.com/distributed_lab/ape"
)

func Run(ctx context.Context, cfg config.Config) {
	log := cfg.Log().WithFields(logan.F{
		"service": "proofs-cleaner",
	})

	proofsCleaner := proofs_cleaner.NewProofsCleaner(log, cfg.Storage(), cfg.LinkConfig(), cfg.RunningPeriodsConfig())
	log.Info("starting proofs-cleaner")
	go proofsCleaner.Run(ctx)

	r := chi.NewRouter()

	const slowRequestDurationThreshold = time.Second
	ape.DefaultMiddlewares(r, cfg.Log(), slowRequestDurationThreshold)

	r.Use(
		ape.CtxMiddleware(
			handlers.CtxLog(cfg.Log()),
			handlers.CtxStorage(cfg.Storage()),
		),
	)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/proofs", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(handlers.AuthMiddleware())
				r.Post("/", handlers.ProofCreate)
			})
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.ProofByID)
			})
			r.Route("/user/{user_did}", func(r chi.Router) {
				r.Use(handlers.AuthMiddleware())
				r.Get("/", handlers.ProofsByUserDID)
			})
		})
	})

	cfg.Log().WithFields(logan.F{
		"service": "api",
		"addr":    cfg.Listener().Addr(),
	}).Info("starting api")

	ape.Serve(ctx, r, cfg, ape.ServeOpts{})
}
