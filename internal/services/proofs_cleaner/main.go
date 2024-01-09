package proofs_cleaner

import (
	"context"
	"github.com/rarimo/rarime-link-svc/internal/config"
	"github.com/rarimo/rarime-link-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"time"
)

type ProofsCleaner struct {
	log                  *logan.Entry
	storage              data.Storage
	cfg                  config.LinkConfig
	proofsCleanerPeriods config.RunningPeriod
}

func NewProofsCleaner(log *logan.Entry, storage data.Storage, cfg config.LinkConfig, periods config.RunningPeriodsConfig) ProofsCleaner {
	return ProofsCleaner{
		log:                  log,
		storage:              storage,
		cfg:                  cfg,
		proofsCleanerPeriods: periods.ProofsCleaner,
	}
}

func (p ProofsCleaner) Run(ctx context.Context) {
	running.WithBackOff(
		ctx, p.log, "proofs-cleaner",
		p.clean,
		p.proofsCleanerPeriods.NormalPeriod,
		p.proofsCleanerPeriods.MinAbnormalPeriod,
		p.proofsCleanerPeriods.MaxAbnormalPeriod,
	)

}

func (p ProofsCleaner) clean(_ context.Context) error {
	proofs, err := p.storage.ProofQ().SelectAllCtx(context.Background())
	if err != nil {
		return err
	}

	for _, proof := range proofs {
		if proof.CreatedAt.Add(p.cfg.MaxExpirationTime).Before(time.Now()) {
			links, err := p.storage.LinksToProofQ().GetLinksToProofsByProofID(context.Background(), proof.ID)
			if err != nil {
				return err
			}

			if len(links) > 0 {
				for _, link := range links {
					err = p.storage.LinksToProofQ().Delete(link)
					if err != nil {
						return err
					}
					err = p.storage.LinkQ().Delete(&data.Link{
						ID: link.LinkID,
					})
					if err != nil {
						return err
					}
				}
			}

			err = p.storage.ProofQ().Delete(proof)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
