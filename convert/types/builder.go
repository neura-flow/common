package types

import (
	"github.com/neura-flow/common/named"
	"github.com/neura-flow/common/util"
)

type Builder interface {
	Build(name named.Name, bundle *Bundle) (named.Named, error)
}

type BuildFunc func(name named.Name, bundle *Bundle) (named.Named, error)

func (f BuildFunc) Build(name named.Name, bundle *Bundle) (named.Named, error) {
	return f(name, bundle)
}

func NewBuilder(bf BuildFunc) Builder {
	return bf
}

type Wrapper struct {
	named.Named
	Target Converter
}

func NewWrapper(c Converter) *Wrapper {
	return &Wrapper{
		Target: c,
	}
}

func DefaultIfBlank(s string) string {
	if util.IsBlank(&s) {
		return "default"
	}
	return s
}
