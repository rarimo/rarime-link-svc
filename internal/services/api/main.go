package api

import (
	"context"
	"time"

	"github.com/go-chi/chi"
	"github.com/rarimo/rarime-link-svc/internal/config"
	"github.com/rarimo/rarime-link-svc/internal/services/api/handlers"
	"github.com/rarimo/rarime-link-svc/internal/services/proofs_cleaner"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
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
			handlers.CtxPoints(cfg.Points()),
		),
	)

	r.Route("/integrations/rarime-link-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/proofs", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(handlers.AuthMiddleware(cfg.Auth(), cfg.Log()))
					r.Get("/", handlers.GetProofs)
					r.Post("/", handlers.CreateProof)
				})

				r.Route("/{id}", func(r chi.Router) {
					r.Use(handlers.AuthMiddleware(cfg.Auth(), cfg.Log()))
					r.Get("/", handlers.ProofByID)
				})
			})

			r.Route("/links", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(handlers.AuthMiddleware(cfg.Auth(), cfg.Log()))
					r.Get("/", handlers.GetLinks)
					r.Post("/", handlers.CreateProofLink)
				})

				r.Route("/{link_id}", func(r chi.Router) {
					r.Use(handlers.OptAuthMiddleware(cfg.Auth(), cfg.Log()))
					r.Get("/", handlers.GetLinkByID)
				})
			})
		})
	})

	cfg.Log().WithFields(logan.F{
		"service": "api",
		"addr":    cfg.Listener().Addr(),
	}).Info("starting api")

	ape.Serve(ctx, r, cfg, ape.ServeOpts{})
}
