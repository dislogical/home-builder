package config

import homebuilder "github.com/dislogical/home-builder/pkg"

var _ homebuilder.ResourceStatusQueryable = (*Config)(nil)

type Config struct{}

func Prepare(resource *homebuilder.Resource) error {
	resource.Impl = &Config{}

	return nil
}

func (p *Config) GetStatus() (homebuilder.ResourceStatus, error) {
	return homebuilder.StatusUpToDate, nil
}

func init() {
	homebuilder.RegisterResourceFactory("config", Prepare)
}
