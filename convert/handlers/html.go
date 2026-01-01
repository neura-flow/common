package handlers

import (
	"context"
	"fmt"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/neura-flow/common/convert/types"
	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/mimetype"
	"github.com/neura-flow/common/named"
	"github.com/neura-flow/common/pandoc"
	"github.com/neura-flow/common/util"
)

func init() {
	bld := types.NewBuilder(func(name named.Name, bundle *types.Bundle) (named.Named, error) {
		c, err := NewHtmlConverter(name.ShortName().Name(), bundle)
		if err != nil {
			return nil, err
		}
		return types.NewWrapper(c), nil
	})
	for _, suffix := range []string{"htm", "html"} {
		for _, engine := range []string{types.EngineHtmlToMarkdown, types.EnginePandoc} {
			types.Register(types.NewComponent(named.Name(suffix+"."+engine), bld))
		}
	}
}

type HtmlConverter struct {
	engine string
	logger log.Logger
	bundle *types.Bundle
}

func NewHtmlConverter(engine string, bundle *types.Bundle) (*HtmlConverter, error) {
	engines := []string{types.EngineHtmlToMarkdown, types.EnginePandoc}
	if !util.OneOf(engine, engines) {
		return nil, fmt.Errorf("engine should be one of [%s]", strings.Join(engines, ", "))
	}
	return &HtmlConverter{
		logger: bundle.Logger,
		engine: engine,
	}, nil
}

func (c *HtmlConverter) Do(_ context.Context, doc *types.Document, insId string, _ map[string]any) ([]byte, error) {
	if mimeType, ok := mimetype.Valid(doc.Content, []string{"text/plain", "text/html", "application/xhtml+xml", "text/xml", "application/xml"}); !ok {
		return nil, fmt.Errorf("mimetype: %s not supported", mimeType)
	}
	if strings.EqualFold(c.engine, types.EngineHtmlToMarkdown) {
		str, err := htmltomarkdown.ConvertString(string(doc.Content))
		if err != nil {
			return nil, err
		}
		return []byte(str), nil
	} else if strings.EqualFold(c.engine, types.EnginePandoc) {
		cli, ok := c.bundle.Pandocs[insId]
		if !ok {
			return nil, fmt.Errorf("pandoc: %s not supported", insId)
		}
		return cli.Convert(doc.Content, pandoc.Options{
			From:    "html",
			To:      "markdown_mmd",
			DataDir: fmt.Sprintf("%s/pandoc", c.bundle.PandocConfigs[insId].SafeDir),
		})
	}
	return doc.Content, nil
}
