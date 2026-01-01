package tongyi_ocr

import (
	"encoding/base64"
	"fmt"

	"github.com/neura-flow/common/util"
)

type ImageMessageContent struct {
	Type      string `json:"type,omitempty"`
	ImageUrl  any    `json:"image_url,omitempty"`
	MinPixels *int   `json:"min_pixels,omitempty"`
	MaxPixels *int   `json:"max_pixels,omitempty"`
	Text      string `json:"text,omitempty"`
}

type ImageTextMessageContent struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

func NewImageUrlMessageContents(url string) []*ImageMessageContent {
	return []*ImageMessageContent{
		{
			Type:      "image_url",
			ImageUrl:  url,
			MinPixels: util.Int(3136),
			MaxPixels: util.Int(1003520),
		},
		{
			Type: "text",
			Text: "Read all the text in the image.",
		},
	}
}

// NewImageBlobMessageContents
// contentType支持列表见: https://help.aliyun.com/zh/model-studio/user-guide/vision?spm=a2c4g.11186623.0.0.5cde453aRLsF0u#da33480805fjh
func NewImageBlobMessageContents(contentType string, image []byte) []*ImageMessageContent {
	imageUrl := fmt.Sprintf("data:%s;base64,%s", contentType, base64.StdEncoding.EncodeToString(image))
	return []*ImageMessageContent{
		{
			Type: "image_url",
			ImageUrl: map[string]string{
				"url": imageUrl,
			},
			MinPixels: util.Int(3136),
			MaxPixels: util.Int(1003520),
		},
		{
			Type: "text",
			Text: "Read all the text in the image.",
		},
	}
}
