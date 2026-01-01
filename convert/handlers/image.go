package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/neura-flow/common/convert/types"
	"github.com/neura-flow/common/idp"
	"github.com/neura-flow/common/llm/tongyi_ocr"
	llmtype "github.com/neura-flow/common/llm/types"
	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/mimetype"
	"github.com/neura-flow/common/named"
	"github.com/neura-flow/common/util"
)

var (
	textInSupportMimeTypes = []string{"image/png", "image/jpeg", "image/bmp", "image/tiff", "image/webp"}
	tongYiSupportMimeTypes = []string{"image/bmp", "image/icns", "image/x-icon", "image/jpeg", "image/jp2", "image/png", "image/sgi", "image/tiff", "image/webp"}
)

func init() {
	bld := types.NewBuilder(func(name named.Name, bundle *types.Bundle) (named.Named, error) {
		c, err := NewImageConverter(name.ShortName().Name(), bundle)
		if err != nil {
			return nil, err
		}
		return types.NewWrapper(c), nil
	})
	for _, suffix := range []string{"png", "jpg", "jpeg", "webp", "bmp", "tiff"} {
		types.Register(types.NewComponent(named.Name(suffix+"."+types.EngineTextIn), bld))
	}
	for _, suffix := range []string{"png", "jpg", "jpeg", "webp", "bmp", "icns", "ico", "jp2", "tiff"} {
		types.Register(types.NewComponent(named.Name(suffix+"."+types.EngineTongYiOcr), bld))
	}
}

type ImageConverter struct {
	engine string
	logger log.Logger
	bundle *types.Bundle
}

func NewImageConverter(engine string, bundle *types.Bundle) (*ImageConverter, error) {
	engines := []string{types.EngineTextIn, types.EngineTongYiOcr}
	if !util.OneOf(engine, engines) {
		return nil, fmt.Errorf("engine should be one of [%s]", strings.Join(engines, ", "))
	}

	c := &ImageConverter{
		engine: engine,
		logger: bundle.Logger,
		bundle: bundle,
	}
	return c, nil
}

func (c *ImageConverter) Do(ctx context.Context, doc *types.Document, insId string, opts map[string]any) ([]byte, error) {
	// png, jpg, jpeg, bmp, tiff, webp,
	if strings.EqualFold(c.engine, types.EngineTextIn) {
		cli, ok := c.bundle.TextIns[insId]
		if !ok {
			return nil, fmt.Errorf("textin: %s not supported", insId)
		}
		if mimeType, ok := mimetype.Valid(doc.Content, textInSupportMimeTypes); !ok {
			return nil, fmt.Errorf("mimeType: %s not support", mimeType)
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
	} else if strings.EqualFold(c.engine, types.EngineTongYiOcr) {
		cli, ok := c.bundle.TongYiOcrs[insId]
		if !ok {
			return nil, fmt.Errorf("tongyiocr: %s not supported", insId)
		}
		if mimeType, ok := mimetype.Valid(doc.Content, tongYiSupportMimeTypes); !ok {
			return nil, fmt.Errorf("mimeType: %s not support", mimeType)
		}
		var model string
		if c.bundle.TongYiConfigs[insId].Model != "" {
			model = c.bundle.TongYiConfigs[insId].Model
		} else {
			model = "qwen-vl-ocr"
		}
		completion := llmtype.NewTongYiCompletionRequest(model)
		completion.Messages = []any{
			llmtype.UserMessage{
				Role: "user", Content: tongyi_ocr.NewImageBlobMessageContents("image/jpeg", doc.Content),
			},
		}
		resp, err := cli.Call(ctx, completion)
		if err != nil {
			return nil, err
		}
		if resp.Error != nil {
			return nil, fmt.Errorf("%v", resp.Error.Error.Message)
		}
		return []byte(resp.Content.Choices[0].Message.Content), nil
	}
	return doc.Content, nil
}
