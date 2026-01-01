package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/neura-flow/common/convert/types"
	"github.com/neura-flow/common/idp"
	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/named"
	"github.com/neura-flow/common/util"
)

func init() {
	bld := types.NewBuilder(func(name named.Name, bundle *types.Bundle) (named.Named, error) {
		c, err := NewDocConverter(name.ShortName().Name(), bundle)
		if err != nil {
			return nil, err
		}
		return types.NewWrapper(c), nil
	})
	for _, suffix := range []string{"doc"} {
		types.Register(types.NewComponent(named.Name(suffix+"."+types.EngineTextIn), bld))
	}
}

type DocConverter struct {
	engine string
	logger log.Logger
	bundle *types.Bundle
}

func NewDocConverter(engine string, bundle *types.Bundle) (*DocConverter, error) {
	engines := []string{types.EngineTextIn}
	if !util.OneOf(engine, engines) {
		return nil, fmt.Errorf("engine should be one of [%s]", strings.Join(engines, ", "))
	}
	c := &DocConverter{
		engine: engine,
		logger: bundle.Logger,
		bundle: bundle,
	}
	return c, nil
}

func (c *DocConverter) Do(ctx context.Context, doc *types.Document, insId string, opts map[string]any) ([]byte, error) {
	if strings.EqualFold(c.engine, types.EngineTextIn) {
		cli, ok := c.bundle.TextIns[insId]
		if !ok {
			return nil, fmt.Errorf("textin instance: %s not found", insId)
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
