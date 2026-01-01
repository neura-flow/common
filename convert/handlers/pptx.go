package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/neura-flow/common/convert/types"
	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/markitdown"
	"github.com/neura-flow/common/mimetype"
	"github.com/neura-flow/common/named"
	"github.com/neura-flow/common/util"
)

func init() {
	bld := types.NewBuilder(func(name named.Name, bundle *types.Bundle) (named.Named, error) {
		c, err := NewPptxConverter(name.ShortName().Name(), bundle)
		if err != nil {
			return nil, err
		}
		return types.NewWrapper(c), nil
	})
	for _, suffix := range []string{"pptx"} {
		for _, engine := range []string{types.EngineMarkItDown} {
			types.Register(types.NewComponent(named.Name(suffix+"."+engine), bld))
		}
	}
}

type PptxConverter struct {
	engine string
	logger log.Logger
	bundle *types.Bundle
}

func NewPptxConverter(engine string, bundle *types.Bundle) (*PptxConverter, error) {
	engines := []string{types.EngineMarkItDown}
	if !util.OneOf(engine, engines) {
		return nil, fmt.Errorf("engine should be one of [%s]", strings.Join(engines, ", "))
	}
	c := &PptxConverter{
		engine: engine,
		logger: bundle.Logger,
		bundle: bundle,
	}
	return c, nil
}

func (c *PptxConverter) Do(ctx context.Context, doc *types.Document, insId string, opts map[string]any) ([]byte, error) {
	if mimeType, ok := mimetype.Valid(doc.Content, []string{"application/vnd.openxmlformats-officedocument.presentationml.presentation"}); !ok {
		return nil, fmt.Errorf("mimetype: %s not supported", mimeType)
	}
	if strings.EqualFold(c.engine, types.EngineMarkItDown) {
		cli, ok := c.bundle.MarkItDowns[insId]
		if !ok {
			return nil, fmt.Errorf("markitdown instance: %s not supported", insId)
		}
		return cli.Convert(doc.Content, markitdown.Options{
			DataDir: fmt.Sprintf("%s/makeitdown", c.bundle.MarkItDownConfigs[insId].SafeDir),
		})
	}
	return doc.Content, nil
}
