package homebuilder

type ResourceFactory func(any) Resource

var resourceFactories map[string]ResourceFactory

func RegisterResourceFactory(name string, factory ResourceFactory) {
	resourceFactories[name] = factory
}
