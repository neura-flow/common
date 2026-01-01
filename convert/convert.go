package convert

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/neura-flow/common/convert/types"
	"github.com/neura-flow/common/idp"
	"github.com/neura-flow/common/llm/tongyi_ocr"
	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/markitdown"
	"github.com/neura-flow/common/mimetype"
	"github.com/neura-flow/common/named"
	"github.com/neura-flow/common/pandoc"
	"github.com/neura-flow/common/util"
)

type Option func(*types.Config)

func AddPandocConfig(cfg *types.PandocConfig) Option {
	return func(c *types.Config) {
		if c.PandocConfigs == nil {
			c.PandocConfigs = make([]*types.PandocConfig, 0)
		}
		cfg.ID = types.DefaultIfBlank(cfg.ID)
		c.PandocConfigs = append(c.PandocConfigs, cfg)
	}
}

func AddMarkItDownConfig(cfg *types.MarkItDownConfig) Option {
	return func(c *types.Config) {
		if c.MarkItDownConfigs == nil {
			c.MarkItDownConfigs = make([]*types.MarkItDownConfig, 0)
		}
		cfg.ID = types.DefaultIfBlank(cfg.ID)
		c.MarkItDownConfigs = append(c.MarkItDownConfigs, cfg)
	}
}

func AddTextInConfig(cfg *types.TextInConfig) Option {
	return func(c *types.Config) {
		if c.TextInConfigs == nil {
			c.TextInConfigs = make([]*types.TextInConfig, 0)
		}
		cfg.ID = types.DefaultIfBlank(cfg.ID)
		c.TextInConfigs = append(c.TextInConfigs, cfg)
	}
}

func AddTongYiOcrConfig(cfg *types.TongYiOcrConfig) Option {
	return func(c *types.Config) {
		if c.TongYiConfigs == nil {
			c.TongYiConfigs = make([]*types.TongYiOcrConfig, 0)
		}
		cfg.ID = types.DefaultIfBlank(cfg.ID)
		c.TongYiConfigs = append(c.TongYiConfigs, cfg)
	}
}

type Converter struct {
	logger log.Logger
	cfg    *types.Config
	bundle *types.Bundle
	m      sync.Map
}

func New(logger log.Logger, opts ...Option) (*Converter, error) {
	c := &Converter{
		logger: logger,
		cfg:    &types.Config{},
	}
	for _, opt := range opts {
		opt(c.cfg)
	}
	var err error
	if c.bundle, err = c.initBundle(); err != nil {
		return nil, err
	}
	cm := types.GlobalManager()
	for _, key := range util.SortString(cm.Keys()) {
		name := named.Name(key)
		if comp, _ := cm.Load(name); comp != nil {
			if conv, err0 := comp.Build(name, c.bundle); err0 != nil {
				logger.Warnf("Failed to build converter: %s err: %v", key, err0)
			} else {
				c.m.Store(key, conv)
				logger.Infof("Converter %s build success", key)
			}
		}
	}
	return c, nil
}

func (c *Converter) initBundle() (*types.Bundle, error) {
	bundle := &types.Bundle{
		Logger:            c.logger,
		PandocConfigs:     make(map[string]*types.PandocConfig),
		MarkItDownConfigs: make(map[string]*types.MarkItDownConfig),
		TextInConfigs:     make(map[string]*types.TextInConfig),
		TongYiConfigs:     make(map[string]*types.TongYiOcrConfig),
		Pandocs:           make(map[string]*pandoc.Client),
		TextIns:           make(map[string]idp.Client),
		TongYiOcrs:        make(map[string]*tongyi_ocr.Client),
		MarkItDowns:       make(map[string]*markitdown.Client),
	}
	for _, item := range c.cfg.PandocConfigs {
		if cli, err := pandoc.New(c.logger, &pandoc.Config{
			Command:    item.Command,
			TimeoutSec: item.TimeoutSec,
			SafeDir:    item.SafeDir,
		}); err != nil {
			return nil, err
		} else {
			id := types.DefaultIfBlank(item.ID)
			bundle.PandocConfigs[id] = item
			bundle.Pandocs[id] = cli
		}
	}
	for _, item := range c.cfg.MarkItDownConfigs {
		if cli, err := markitdown.New(c.logger, &markitdown.Config{
			Command:    item.Command,
			TimeoutSec: item.TimeoutSec,
			SafeDir:    item.SafeDir,
		}); err != nil {
			return nil, err
		} else {
			id := types.DefaultIfBlank(item.ID)
			bundle.MarkItDownConfigs[id] = item
			bundle.MarkItDowns[id] = cli
		}
	}
	for _, item := range c.cfg.TextInConfigs {
		cli := idp.NewTextInOcr(&idp.TextInConfig{
			AppID:     item.AppID,
			AppSecret: item.AppSecret,
			Host:      item.Host,
		})
		id := types.DefaultIfBlank(item.ID)
		bundle.TextInConfigs[id] = item
		bundle.TextIns[id] = cli
	}
	for _, item := range c.cfg.TongYiConfigs {
		cli := tongyi_ocr.New(&tongyi_ocr.Config{
			ServerUrl: item.ServerUrl,
			Secret:    item.Secret,
		})

		id := types.DefaultIfBlank(item.ID)
		bundle.TongYiConfigs[id] = item
		bundle.TongYiOcrs[id] = cli
	}
	return bundle, nil
}

// Do 把指定的文档转化为 markdown, 可指定 engine 和 insId , 如未指定，那么从默认列表中读取
func (c *Converter) Do(ctx context.Context, doc *types.Document, engine, insId string, opts map[string]any) ([]byte, error) {
	suffix := c.detectSuffix(doc)
	target, err := c.lookEngine(suffix, engine)
	if err != nil {
		return nil, err
	}

	ins, ok := c.m.Load(fmt.Sprintf("%s.%s", suffix, target))
	if !ok {
		return nil, fmt.Errorf("suffix: %s 's engine not found", suffix)
	}

	start := time.Now()
	resp, err := ins.(*types.Wrapper).Target.(types.Converter).Do(ctx, doc, types.DefaultIfBlank(insId), opts)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start).Milliseconds()
	c.logger.Infof("Convert from: %s size: %d engine: %s insId: %s duration: %d ms", suffix, len(doc.Content), target, insId, duration)

	return resp, nil
}

func (c *Converter) detectSuffix(doc *types.Document) string {
	// 通过 mimeType 自动识别 suffix
	var suffix string
	if v, ok := mimetype.GetSuffixes()[mimetype.Detect(doc.Content)]; ok && !util.OneOf(doc.Suffix, v) {
		suffix = v[0]
	} else {
		suffix = doc.Suffix
	}

	return suffix
}

func (c *Converter) lookEngine(suffix, engine string) (string, error) {
	engines, ok := c.SuffixEngines()[suffix]
	if !ok {
		return "", fmt.Errorf("suffix: %s not support", suffix)
	}
	if len(engines) == 0 {
		return "", fmt.Errorf("suffix: %s did not have engines", suffix)
	}
	if engine != "" && !util.OneOf(engine, engines) {
		return "", fmt.Errorf("engine %s not support for %s", engine, suffix)
	}

	var engine0 string
	if engine == "" {
		engine0 = util.SortString(engines)[0]
	} else {
		engine0 = engine
	}
	return engine0, nil
}

func (c *Converter) SuffixEngines() map[string][]string {
	kvm := make(map[string][]string)
	c.m.Range(func(key, value interface{}) bool {
		arr := strings.Split(key.(string), ".")
		if val, ok := kvm[arr[0]]; ok {
			kvm[arr[0]] = append(val, arr[1])
		} else {
			kvm[arr[0]] = []string{arr[1]}
		}
		return true
	})
	return kvm
}
