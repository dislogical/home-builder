package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/urfave/cli/v3"

	homebuilder "github.com/dislogical/home-builder/pkg"
)

var init_ *cli.Command = &cli.Command{
	Name: "init",
	Action: func(ctx context.Context, c *cli.Command) error {
		configDir := filepath.Join(xdg.ConfigHome, "home-builder")
		configPath := filepath.Join(configDir, "config.cue")

		_, err := os.Stat(configPath)
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("configuration file already exists: %s", configPath)
		}

		err = os.MkdirAll(configDir, 0o750)
		if err != nil {
			return err //nolint:wrapcheck
		}

		configF, err := os.Create(configPath)
		if err != nil {
			return err //nolint:wrapcheck
		}

		_, err = configF.WriteString(homebuilder.DefaultConfig)
		if err != nil {
			return err //nolint:wrapcheck
		}

		return nil
	},
}
