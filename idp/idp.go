package idp

import "context"

type Client interface {
	Convert(ctx context.Context, doc *Document) (*Result, error)
}

type Document struct {
	ContentType string
	Content     []byte
	Options     map[string]interface{}
}

type Result struct {
	Content []byte
}
