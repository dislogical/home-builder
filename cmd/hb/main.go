package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"

	"github.com/adrg/xdg"
	"github.com/urfave/cli/v3"

	homebuilder "github.com/dislogical/home-builder/pkg"
	_ "github.com/dislogical/home-builder/pkg/config"
	_ "github.com/dislogical/home-builder/pkg/packages"
)

const Schema = `
home: [_type=string]: [_name=string]: {
	$meta: {
		name: string | *_name
		type: string | *_type
	}
}
`

//nolint:gocognit
func main() {
	cmd := &cli.Command{
		Name: "hb",
		Commands: []*cli.Command{
			diff,
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			cuectx := cuecontext.New()

			insts := load.Instances([]string{"."}, &load.Config{
				Dir: filepath.Join(xdg.ConfigHome, "home-builder"),
			})
			if len(insts) != 1 {
				return errors.New("INTERNAL ERROR: expected exactly 1 instance")
			}
			inst := insts[0]

			schema := cuectx.CompileString(Schema)
			if err := schema.Err(); err != nil {
				return fmt.Errorf("INTERNAL ERROR: compiling schema: %w", err)
			}

			val := cuectx.BuildInstance(inst)
			if err := val.Err(); err != nil {
				return fmt.Errorf("building config: %w", err)
			}

			val = val.Unify(schema)
			if err := val.Err(); err != nil {
				return fmt.Errorf("conforming config to schema: %w", err)
			}

			// Get the "home" object for configuration
			val = val.LookupPath(cue.ParsePath("home"))
			if err := val.Err(); err != nil {
				return fmt.Errorf("finding object \"home\": %w", err)
			}

			modules, err := val.Fields()
			if err != nil {
				return fmt.Errorf("retreiving fields: %w", err)
			}

			for modules.Next() {
				module := modules.Value()
				if err = module.Err(); err != nil {
					return fmt.Errorf("getting module: %w", err)
				}

				if !modules.Selector().IsString() {
					return errors.New("module selector must be a string")
				}

				configs, err := module.Fields()
				if err != nil {
					return fmt.Errorf("retreiving fields: %w", err)
				}

				for configs.Next() {
					config := configs.Value()
					if err = config.Err(); err != nil {
						return fmt.Errorf("getting config: %w", err)
					}

					resource := homebuilder.Resource{
						Config: config,
					}
					err = config.Decode(&resource)
					if err != nil {
						return fmt.Errorf("decoding: %w", err)
					}

					err = resource.Prepare()
					if err != nil {
						return fmt.Errorf("preparing resource: %w", err)
					}

					status := homebuilder.StatusUnknown
					if statusable, ok := resource.Impl.(homebuilder.ResourceStatusQueryable); ok {
						status, err = statusable.GetStatus()
						if err != nil {
							return fmt.Errorf("retrieving status for %s: %w", resource.Meta, err)
						}
					}

					log.Printf("%s: %s", resource.Meta.String(), status)
				}
			}

			return nil
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
