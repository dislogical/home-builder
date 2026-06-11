package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/urfave/cli/v3"

	tea "charm.land/bubbletea/v2"

	homebuilder "github.com/dislogical/home-builder/pkg"
	"github.com/dislogical/home-builder/pkg/bubbles"
	_ "github.com/dislogical/home-builder/pkg/config"
	_ "github.com/dislogical/home-builder/pkg/packages"
)

var main_ = &cli.Command{
	Name: "hb",
	Commands: []*cli.Command{
		init_,
		diff,
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		hbctx := homebuilder.NewContext()

		resources, err := hbctx.Load(filepath.Join(xdg.ConfigHome, "home-builder"))
		if errors.Is(err, homebuilder.ErrConfigNotExist) {
			return errors.New("config directory does not exist, run 'init' to create a default config")
		}
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		program := tea.NewProgram(
			&bubbles.ContextModel{HbCtx: &hbctx, Resources: resources},
			tea.WithContext(ctx),
		)

		_, err = program.Run()

		return err //nolint:wrapcheck
	},
}

func main() {
	err := main_.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
