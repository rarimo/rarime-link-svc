package config

type SubServices interface {
	SubServicesConfig() SubServicesConfig
	SetConfig(config SubServicesConfig)
}

type SubServicesConfig struct {
	ProofsCleaner bool
}

type subServices struct {
	config SubServicesConfig
}

func NewSubServices() SubServices {
	return &subServices{
		config: SubServicesConfig{},
	}
}

func (l *subServices) SubServicesConfig() SubServicesConfig {
	return l.config
}

func (l *subServices) SetConfig(config SubServicesConfig) {
	l.config = config
}
