package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/urfave/cli/v3"

	homebuilder "github.com/dislogical/home-builder/pkg"
)

var diff *cli.Command = &cli.Command{
	Name: "diff",
	Action: func(ctx context.Context, c *cli.Command) error {
		hbctx := homebuilder.NewContext()
		dmp := diffmatchpatch.New()

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

			if diffable, ok := resource.Backend.(homebuilder.ResourceDiffable); ok {
				if status == homebuilder.StatusNeedsUpdate {
					expected, current, err := diffable.GetDiff()
					if err != nil {
						log.Printf("%s: error diffing: %s", resource.Meta, err)
					} else {
						log.Println(resource.Meta.String())
						diffs := dmp.DiffMain(current, expected, true)
						log.Println(dmp.DiffPrettyText(diffs))
					}
				}
			}
		}

		return nil
	},
}
