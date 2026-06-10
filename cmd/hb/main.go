package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/urfave/cli/v3"

	homebuilder "github.com/dislogical/home-builder/pkg"
	_ "github.com/dislogical/home-builder/pkg/config"
	_ "github.com/dislogical/home-builder/pkg/packages"
)

func main() {
	cmd := &cli.Command{
		Name: "hb",
		Commands: []*cli.Command{
			diff,
			test,
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			hbctx := homebuilder.NewContext()

			resources, err := hbctx.Load(filepath.Join(xdg.ConfigHome, "home-builder"))
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			for _, resource := range resources {
				err = resource.Prepare()
				if err != nil {
					return fmt.Errorf("preparing resource: %w", err)
				}

				status, err := resource.Backend.GetStatus()
				if err != nil {
					return fmt.Errorf("retrieving status for %s: %w", resource.Meta, err)
				}

				log.Printf("%s: %s", resource.Meta.String(), status)
			}

			return nil
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
