package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/neura-flow/common/convert/types"
	"github.com/neura-flow/common/idp"
	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/markitdown"
	"github.com/neura-flow/common/mimetype"
	"github.com/neura-flow/common/named"
	"github.com/neura-flow/common/pandoc"
	"github.com/neura-flow/common/util"
)

func init() {
	bld := types.NewBuilder(func(name named.Name, bundle *types.Bundle) (named.Named, error) {
		c, err := NewPDFConverter(name.ShortName().Name(), bundle)
		if err != nil {
			return nil, err
		}
		return types.NewWrapper(c), nil
	})
	for _, suffix := range []string{"pdf"} {
		for _, engine := range []string{types.EngineMarkItDown, types.EnginePandoc, types.EngineTextIn} {
			types.Register(types.NewComponent(named.Name(suffix+"."+engine), bld))
		}
	}
}

type PDFConverter struct {
	engine string
	logger log.Logger
	bundle *types.Bundle
}

func NewPDFConverter(engine string, bundle *types.Bundle) (*PDFConverter, error) {
	engines := []string{types.EngineMarkItDown, types.EnginePandoc, types.EngineTextIn}
	if !util.OneOf(engine, engines) {
		return nil, fmt.Errorf("engine should be one of [%s]", strings.Join(engines, ", "))
	}
	c := &PDFConverter{
		engine: engine,
		logger: bundle.Logger,
		bundle: bundle,
	}
	return c, nil
}

func (c *PDFConverter) Do(ctx context.Context, doc *types.Document, insId string, opts map[string]any) ([]byte, error) {
	if mimeType, ok := mimetype.Valid(doc.Content, []string{"application/pdf"}); !ok {
		return nil, fmt.Errorf("mimetype: %s not supported", mimeType)
	}
	if strings.EqualFold(c.engine, types.EngineMarkItDown) {
		cli, ok := c.bundle.MarkItDowns[insId]
		if !ok {
			return nil, fmt.Errorf("markitdown: %s not supported", insId)
		}
		return cli.Convert(doc.Content, markitdown.Options{
			DataDir: fmt.Sprintf("%s/makeitdown", c.bundle.MarkItDownConfigs[insId].SafeDir),
		})
	} else if strings.EqualFold(c.engine, types.EnginePandoc) {
		cli, ok := c.bundle.Pandocs[insId]
		if !ok {
			return nil, fmt.Errorf("pandoc: %s not supported", insId)
		}
		return cli.Convert(doc.Content, pandoc.Options{
			From:    "pdf",
			To:      "markdown_mmd",
			DataDir: fmt.Sprintf("%s/pandoc", c.bundle.PandocConfigs[insId].SafeDir),
		})
	} else if strings.EqualFold(c.engine, types.EngineTextIn) {
		cli, ok := c.bundle.TextIns[insId]
		if !ok {
			return nil, fmt.Errorf("textin: %s not supported", insId)
		}
		options := idp.Options{
			PageStart:   util.Int(0),
			PageCount:   util.Int(1000),
			TableFlavor: util.String("md"),
			ParseMode:   util.String("scan"),
			Dpi:         util.Int(144),
			PageDetails: util.Int(0),
		}
		kvm := make(map[string]interface{})
		_ = json.Unmarshal(util.ToJsonBytes(options), &kvm)

		doc0 := &idp.Document{
			ContentType: "application/octet-stream",
			Content:     doc.Content,
			Options:     kvm,
		}
		resp, err := cli.Convert(ctx, doc0)
		if err != nil {
			return nil, err
		}
		return resp.Content, nil
	}
	return doc.Content, nil
}
