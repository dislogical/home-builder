package homebuilder

type Resource struct {
	Meta string `json:"$meta" yaml:"$meta"`
	Impl any    `json:"-"     yaml:"-"`
}

type ResourceStatus int

const (
	StatusUpToDate ResourceStatus = iota
	StatusMissing
	StatusNeedsUpdate
)

type ResourceStatusQueryable interface {
	GetStatus() (ResourceStatus, error)
}

type ResourceDiffable interface {
	GetDiff() (string, string, error)
}
