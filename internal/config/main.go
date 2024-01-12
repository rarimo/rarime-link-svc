package config

import (
	"github.com/rarimo/rarime-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-link-svc/internal/data"
	"github.com/rarimo/rarime-link-svc/internal/data/pg"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	comfig.Listenerer
	pgdb.Databaser
	auth.Auther
	LinkConfiger
	RunningPeriods
	SubServices
	Storage() data.Storage
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	auth.Auther
	LinkConfiger
	RunningPeriods
	SubServices

	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:         getter,
		Listenerer:     comfig.NewListenerer(getter),
		Logger:         comfig.NewLogger(getter, comfig.LoggerOpts{}),
		LinkConfiger:   NewLinkConfiger(getter),
		RunningPeriods: NewRunningPeriods(getter),
		SubServices:    NewSubServices(),
		Databaser:      pgdb.NewDatabaser(getter),
		Auther:         auth.NewAuther(getter),
	}
}

func (c *config) Link() LinkConfig {
	return c.LinkConfiger.LinkConfig()
}

func (c *config) RunningPeriod() RunningPeriodsConfig {
	return c.RunningPeriods.RunningPeriodsConfig()
}

func (c *config) SubService() SubServicesConfig {
	return c.SubServices.SubServicesConfig()
}

func (c *config) Storage() data.Storage {
	return pg.New(c.DB().Clone())
}
