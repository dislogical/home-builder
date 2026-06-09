package homebuilder

import (
	"errors"
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

type Context struct {
	Cue *cue.Context

	configSchema cue.Value
}

func NewContext() Context {
	cuectx := cuecontext.New()

	configSchema := cuectx.CompileString(resourceSchema)
	if err := configSchema.Err(); err != nil {
		panic(err)
	}

	return Context{
		Cue:          cuectx,
		configSchema: configSchema,
	}
}

func (ctx *Context) Load(dir string) ([]Resource, error) {
	var config cue.Value

	insts := load.Instances([]string{"."}, &load.Config{
		Dir: dir,
	})
	if len(insts) != 1 {
		return nil, errors.New("INTERNAL ERROR: expected exactly 1 instance")
	}
	inst := insts[0]

	config = ctx.Cue.BuildInstance(inst)
	if err := config.Err(); err != nil {
		return nil, fmt.Errorf("building config: %w", err)
	}

	config = config.Unify(ctx.configSchema)
	if err := config.Err(); err != nil {
		return nil, fmt.Errorf("conforming config to schema: %w", err)
	}

	// Get the "home" object for configuration
	config = config.LookupPath(cue.ParsePath("home"))
	if err := config.Err(); err != nil {
		return nil, fmt.Errorf("finding object \"home\": %w", err)
	}

	modules, err := config.Fields()
	if err != nil {
		return nil, fmt.Errorf("retreiving fields: %w", err)
	}

	resources := []Resource{}

	for modules.Next() {
		module := modules.Value()
		if err = module.Err(); err != nil {
			return nil, fmt.Errorf("getting module: %w", err)
		}

		if !modules.Selector().IsString() {
			return nil, errors.New("module selector must be a string")
		}

		configs, err := module.Fields()
		if err != nil {
			return nil, fmt.Errorf("retreiving fields: %w", err)
		}

		for configs.Next() {
			config := configs.Value()
			if err = config.Err(); err != nil {
				return nil, fmt.Errorf("getting config: %w", err)
			}

			resource := Resource{
				Config: config,
			}
			err = config.Decode(&resource)
			if err != nil {
				return nil, fmt.Errorf("decoding: %w", err)
			}

			resources = append(resources, resource)
		}
	}

	return resources, nil
}
