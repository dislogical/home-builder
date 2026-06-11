package packages

import homebuilder "github.com/dislogical/home-builder/pkg"

var _ homebuilder.ResourceBackend = (*Package)(nil)

type Package struct{}

type Factory struct{}

func (f *Factory) InitBackend(resource *homebuilder.Resource) error {
	resource.Backend = &Package{}

	return nil
}

func (p *Package) GetStatus() (homebuilder.ResourceStatus, error) {
	return homebuilder.StatusUpToDate, nil
}

func init() {
	homebuilder.RegisterResourceFactory("package", &Factory{})
}
