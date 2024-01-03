package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/tokend/keypair/figurekeypair"
	"time"
)

type LinkConfiger interface {
	LinkConfig() LinkConfig
}

func NewLinkConfiger(getter kv.Getter) LinkConfiger {
	return &Link{
		getter: getter,
	}
}

type Link struct {
	LinkOnce comfig.Once
	getter   kv.Getter
}

type LinkConfig struct {
	MaxExpirationTime time.Duration `fig:"max_expiration_time,required"`
}

func (e *Link) LinkConfig() LinkConfig {
	return e.LinkOnce.Do(func() interface{} {
		var result LinkConfig

		err := figure.
			Out(&result).
			With(figure.BaseHooks, figurekeypair.Hooks).
			From(kv.MustGetStringMap(e.getter, "link")).
			Please()
		if err != nil {
			panic(err)
		}

		return result
	}).(LinkConfig)
}
