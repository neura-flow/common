package tongyi_file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type Config struct {
	ServerUrl string `json:"serverUrl"`
	Secret    string `json:"secret"`
}

type Client struct {
	cfg *Config
}

func New(cfg *Config) (*Client, error) {
	if cfg.ServerUrl == "" {
		cfg.ServerUrl = "https://dashscope.aliyuncs.com/compatible-mode/v1/files"
	}
	return &Client{cfg: cfg}, nil
}

func (c *Client) Create(filename string, params map[string]string, data []byte) (*Response, error) {
	body := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(body)

	// file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	n, err := part.Write(data)
	if err != nil {
		return nil, err
	}
	if n != len(data) {
		return nil, fmt.Errorf("failed to copy file to multipart writer")
	}
	// params
	for k, v := range params {
		if err = writer.WriteField(k, v); err != nil {
			return nil, err
		}
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}
	defer writer.Close()

	req, err := http.NewRequest(http.MethodPost, c.cfg.ServerUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.Secret)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var fileInfo *FileInfo
	if err = json.Unmarshal(resBody, &fileInfo); err != nil {
		return nil, err
	}
	if fileInfo != nil && fileInfo.Id != "" {
		return &Response{
			File: fileInfo,
		}, nil
	}

	var errorInfo *ErrorInfo
	if err = json.Unmarshal(resBody, &errorInfo); err != nil {
		return nil, err
	}
	return &Response{
		Error: errorInfo,
	}, nil
}

type Response struct {
	File  *FileInfo  `json:"file"`
	Error *ErrorInfo `json:"error"`
}

type FileInfo struct {
	Id        string `json:"id"`
	Bytes     int64  `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Object    string `json:"object"`
	Purpose   string `json:"purpose"`
	Status    string `json:"status"`
}

type ErrorInfo struct {
	RequestId string              `json:"requestId"`
	Error     *ErrorResponseError `json:"error"`
}

type ErrorResponseError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   any    `json:"param"`
	Code    any    `json:"code"`
}
