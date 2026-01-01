package types

import (
	"context"

	"github.com/neura-flow/common/idp"
	"github.com/neura-flow/common/llm/tongyi_ocr"
	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/markitdown"
	"github.com/neura-flow/common/pandoc"
)

const (
	EngineMarkItDown     = "markitdown"
	EnginePandoc         = "pandoc"
	EngineHtmlToMarkdown = "htmltomarkdown"
	EngineTextIn         = "textin"
	EngineTongYiOcr      = "tongyi_ocr"
	NA                   = "n/a"
	OctetStream          = "application/octet-stream"
	Default              = "default"
)

type Document struct {
	Suffix  string
	Content []byte
}

type Converter interface {
	Do(ctx context.Context, doc *Document, insId string, params map[string]any) ([]byte, error)
}

type PandocConfig struct {
	ID         string `json:"id"`
	Command    string `json:"command"`
	TimeoutSec int    `json:"timeoutSec"`
	SafeDir    string `json:"safeDir"`
}

type MarkItDownConfig struct {
	ID         string `json:"id"`
	Command    string `json:"command"`
	TimeoutSec int    `json:"timeoutSec"`
	SafeDir    string `json:"safeDir"`
}

type TextInConfig struct {
	ID        string `json:"id"`
	AppID     string `json:"appId"`
	AppSecret string `json:"appSecret"`
	Host      string `json:"host"`
}

type TongYiOcrConfig struct {
	ID        string `json:"id"`
	ServerUrl string `json:"serverUrl"`
	Secret    string `json:"secret"`
	Model     string `json:"model"`
}

type Config struct {
	PandocConfigs     []*PandocConfig
	MarkItDownConfigs []*MarkItDownConfig
	TextInConfigs     []*TextInConfig
	TongYiConfigs     []*TongYiOcrConfig
}

type Bundle struct {
	Logger            log.Logger
	PandocConfigs     map[string]*PandocConfig
	MarkItDownConfigs map[string]*MarkItDownConfig
	TextInConfigs     map[string]*TextInConfig
	TongYiConfigs     map[string]*TongYiOcrConfig
	MarkItDowns       map[string]*markitdown.Client
	Pandocs           map[string]*pandoc.Client
	TextIns           map[string]idp.Client
	TongYiOcrs        map[string]*tongyi_ocr.Client
}
