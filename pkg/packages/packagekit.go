package packages

import homebuilder "github.com/dislogical/home-builder/pkg"

var _ homebuilder.ResourceStatusQueryable = (*Package)(nil)

type Package struct{}

func New(config any) homebuilder.Resource {
	return homebuilder.Resource{
		Impl: Package{},
	}
}

func (p *Package) GetStatus() (homebuilder.ResourceStatus, error) {
	return homebuilder.StatusUpToDate, nil
}

func init() {
	homebuilder.RegisterResourceFactory("package", New)
}
