package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"cuelang.org/go/pkg/encoding/json"
	"cuelang.org/go/pkg/encoding/toml"
	"cuelang.org/go/pkg/encoding/yaml"

	"github.com/adrg/xdg"

	homebuilder "github.com/dislogical/home-builder/pkg"
)

var (
	_ homebuilder.ResourceBackend  = (*Config)(nil)
	_ homebuilder.ResourceDiffable = (*Config)(nil)
)

type Config struct {
	FileName        string `json:"$filename"`
	Encoding        string `json:"$encoding"`
	Content         string `json:"-"`
	ExistingContent string `json:"-"`
}

func Prepare(resource *homebuilder.Resource) error {
	config := &Config{}

	err := resource.Config.Decode(config)
	if err != nil {
		return fmt.Errorf("decoding config file: %w", err)
	}

	// Encode the content
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

	// Translate the file name
	tmp, err := template.New("filename").Parse(config.FileName)
	if err != nil {
		return fmt.Errorf("parsing filename \"%s\": %w", config.FileName, err)
	}
	var filenameBuilder strings.Builder
	err = tmp.Execute(&filenameBuilder, map[string]any{
		"xdg": map[string]any{
			"config": xdg.ConfigHome,
		},
	})
	if err != nil {
		return fmt.Errorf("executing filename template \"%s\": %w", config.FileName, err)
	}
	config.FileName = filenameBuilder.String()

	resource.Backend = config

	return nil
}

// GetStatus implements [homebuilder.ResourceBackend].
func (p *Config) GetStatus() (homebuilder.ResourceStatus, error) {
	contents, err := os.ReadFile(p.FileName)
	if errors.Is(err, os.ErrNotExist) {
		return homebuilder.StatusMissing, nil
	}

	if err != nil {
		return homebuilder.StatusUnknown, fmt.Errorf("reading file \"%s\": %w", p.FileName, err)
	}

	p.ExistingContent = string(contents)

	if p.ExistingContent != p.Content {
		return homebuilder.StatusNeedsUpdate, nil
	}

	return homebuilder.StatusUpToDate, nil
}

// GetDiff implements [homebuilder.ResourceDiffable].
func (p *Config) GetDiff() (string, string, error) {
	return p.Content, p.ExistingContent, nil
}

func init() {
	homebuilder.RegisterResourceFactory("config", Prepare)
}
