package config

import (
	"time"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/keypair/figurekeypair"
)

type RunningPeriods interface {
	RunningPeriodsConfig() RunningPeriodsConfig
}

type RunningPeriod struct {
	NormalPeriod      time.Duration `fig:"normal_period,required"`
	MinAbnormalPeriod time.Duration `fig:"min_abnormal_period,required"`
	MaxAbnormalPeriod time.Duration `fig:"max_abnormal_period,required"`
}

type RunningPeriodsConfig struct {
	ProofsCleaner RunningPeriod `fig:"proofs_cleaner,required"`
}

type runningPeriods struct {
	once   comfig.Once
	getter kv.Getter
}

func NewRunningPeriods(getter kv.Getter) RunningPeriods {
	return &runningPeriods{
		getter: getter,
	}
}

func (r *runningPeriods) RunningPeriodsConfig() RunningPeriodsConfig {
	return r.once.Do(func() interface{} {
		var cfg RunningPeriodsConfig
		err := figure.
			Out(&cfg).
			With(figure.BaseHooks, figurekeypair.Hooks).
			From(kv.MustGetStringMap(r.getter, "running_periods")).
			Please()

		if err != nil {
			panic(errors.Wrap(err, "failed to figure out running periods"))
		}
		return cfg
	}).(RunningPeriodsConfig)
}
