package idp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type TextInConfig struct {
	AppID     string `json:"appId"`
	AppSecret string `json:"appSecret"`
	Host      string `json:"host"`
}

type TextInOcr struct {
	cfg *TextInConfig
}

func NewTextInOcr(cfg *TextInConfig) *TextInOcr {
	return &TextInOcr{cfg: cfg}
}

func (ti *TextInOcr) Convert(ctx context.Context, doc *Document) (*Result, error) {
	serverUrl := ti.serverUrl()
	req, err := http.NewRequest("POST", serverUrl, bytes.NewBuffer(doc.Content))
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-ti-app-id", ti.cfg.AppID)
	req.Header.Set("x-ti-secret-code", ti.cfg.AppSecret)
	var contentType string
	if doc.ContentType == "" {
		contentType = "application/octet-stream"
	} else {
		contentType = doc.ContentType
	}
	req.Header.Set("Content-Type", contentType)

	var q = url.Values{}
	for k, v := range doc.Options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jsonData TextInHttpResponse
	if err = json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
		fmt.Println("Error decoding response:", err)
		return nil, err
	}

	if jsonData.Code != http.StatusOK {
		return nil, errors.New(jsonData.Message)
	}

	return &Result{
		Content: []byte(jsonData.Result.Markdown),
	}, nil
}

func (ti *TextInOcr) serverUrl() string {
	return fmt.Sprintf("%s/ai/service/v1/pdf_to_markdown", ti.cfg.Host)
}

type TextInHttpResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		Markdown string `json:"markdown"`
	} `json:"result"`
}

type Options struct {
	PdfPwd            *string `json:"pdf_pwd,omitempty"`
	Dpi               *int    `json:"dpi,omitempty"`
	PageStart         *int    `json:"page_start,omitempty"`
	PageCount         *int    `json:"page_count,omitempty"`
	ApplyDocumentTree *int    `json:"apply_document_tree,omitempty"`
	MarkdownDetails   *int    `json:"markdown_details,omitempty"`
	TableFlavor       *string `json:"table_flavor,omitempty"`
	GetImage          *string `json:"get_image,omitempty"`
	ParseMode         *string `json:"parse_mode,omitempty"`
	PageDetails       *int    `json:"page_details,omitempty"`
}

func getFileContent(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}
