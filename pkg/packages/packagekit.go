package packages

import homebuilder "github.com/dislogical/home-builder/pkg"

var _ homebuilder.ResourceStatusQueryable = (*Package)(nil)

type Package struct{}

func Prepare(resource *homebuilder.Resource) error {
	resource.Impl = &Package{}

	return nil
}

func (p *Package) GetStatus() (homebuilder.ResourceStatus, error) {
	return homebuilder.StatusUpToDate, nil
}

func init() {
	homebuilder.RegisterResourceFactory("package", Prepare)
}
