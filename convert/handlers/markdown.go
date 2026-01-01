package handlers

import (
	"context"

	"github.com/neura-flow/common/convert/types"
	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/named"
)

func init() {
	bld := types.NewBuilder(func(name named.Name, bundle *types.Bundle) (named.Named, error) {
		c, err := NewMarkdownConverter(name.ShortName().Name(), bundle)
		if err != nil {
			return nil, err
		}
		return types.NewWrapper(c), nil
	})
	for _, suffix := range []string{"md"} {
		types.Register(types.NewComponent(named.Name(suffix+".default"), bld))
	}
}

type MarkdownConverter struct {
	engine string
	logger log.Logger
	bundle *types.Bundle
}

func NewMarkdownConverter(engine string, bundle *types.Bundle) (*MarkdownConverter, error) {
	c := &MarkdownConverter{
		engine: engine,
		logger: bundle.Logger,
		bundle: bundle,
	}
	return c, nil
}

func (c *MarkdownConverter) Do(ctx context.Context, doc *types.Document, _ string, opts map[string]any) ([]byte, error) {
	// do nothing
	return doc.Content, nil
}
