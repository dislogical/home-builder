package homebuilder

import (
	"fmt"

	"cuelang.org/go/cue"

	_ "embed"
)

type Metadata struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (m Metadata) String() string {
	return fmt.Sprintf("%s:%s", m.Type, m.Name)
}

type Resource struct {
	Context *Context `json:"-"`

	Meta    Metadata        `json:"$meta"`
	Config  cue.Value       `json:"-"`
	Backend ResourceBackend `json:"-"`
}

//go:embed schema.cue
var resourceSchema string

type ResourceStatus int

const (
	StatusUpToDate ResourceStatus = iota
	StatusMissing
	StatusNeedsUpdate
	StatusUnknown
)

func (rs ResourceStatus) String() string {
	switch rs {
	case StatusUpToDate:
		return "up to date"
	case StatusMissing:
		return "missing"
	case StatusNeedsUpdate:
		return "needs update"

	case StatusUnknown:
		fallthrough
	default:
		return "<unknown>"
	}
}

type ResourceBackend interface {
	GetStatus() (ResourceStatus, error)
}

type ResourceDiffable interface {
	GetDiff() (string, string, error)
}

func (r *Resource) Prepare() error {
	factory, found := resourceFactories[r.Meta.Type]
	if !found {
		return fmt.Errorf("%w: %s", ErrFactoryNotFound, r.Meta)
	}

	err := factory(r)
	if err != nil {
		return err
	}

	if r.Backend == nil {
		return fmt.Errorf("%w: %s", ErrFactoryNooped, r.Meta)
	}

	return nil
}
