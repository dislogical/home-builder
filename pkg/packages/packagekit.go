package packages

import (
	homebuilder "github.com/dislogical/home-builder/pkg"
)

//go:generate ./packagekit/gen.sh

var _ homebuilder.ResourceBackend = (*Package)(nil)

type Package struct {
	Name string `json:"-"`
}

type Factory struct{}

func (f *Factory) InitBackend(resource *homebuilder.Resource) error {
	resource.Backend = &Package{
		Name: resource.Meta.Name,
	}

	return nil
}

func (p *Package) GetStatus() (homebuilder.ResourceStatus, error) {
	return homebuilder.StatusUpToDate, nil
}

func init() {
	homebuilder.RegisterResourceFactory("package", &Factory{})
}
