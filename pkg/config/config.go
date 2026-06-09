package config

import (
	"fmt"

	"cuelang.org/go/pkg/encoding/json"
	"cuelang.org/go/pkg/encoding/toml"
	"cuelang.org/go/pkg/encoding/yaml"

	homebuilder "github.com/dislogical/home-builder/pkg"
)

var _ homebuilder.ResourceBackend = (*Config)(nil)

type Config struct {
	FileName string `json:"$filename"`
	Encoding string `json:"$encoding"`
	Content  string `json:"-"`
}

func Prepare(resource *homebuilder.Resource) error {
	config := &Config{}

	err := resource.Config.Decode(config)
	if err != nil {
		return fmt.Errorf("decoding config file: %w", err)
	}

	if config.Content == "" {
		switch config.Encoding {
		case "json":
			config.Content, err = json.Marshal(resource.Config)
		case "yaml":
			config.Content, err = yaml.Marshal(resource.Config)
		case "toml":
			config.Content, err = toml.Marshal(resource.Config)

		default:
			err = fmt.Errorf("unknown encoding: %s", config.Encoding)
		}
	}
	if err != nil {
		return fmt.Errorf("decoding content: %w", err)
	}

	resource.Backend = config

	return nil
}

func (p *Config) GetStatus() (homebuilder.ResourceStatus, error) {
	return homebuilder.StatusUpToDate, nil
}

func init() {
	homebuilder.RegisterResourceFactory("config", Prepare)
}
