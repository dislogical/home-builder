package homebuilder

import "errors"

type ResourceFactory func(resource *Resource) error

var resourceFactories map[string]ResourceFactory = make(map[string]ResourceFactory)

func RegisterResourceFactory(name string, factory ResourceFactory) {
	resourceFactories[name] = factory
}

var (
	ErrFactoryNotFound = errors.New("resource factory not found")
	ErrFactoryNooped   = errors.New("factory did not initialize resource")
)
