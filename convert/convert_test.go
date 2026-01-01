package convert

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/neura-flow/common/convert/handlers"
	"github.com/neura-flow/common/convert/types"
	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/util"
)

func TestGetSuffixEngines(t *testing.T) {
	fmt.Printf("%s\n", util.ToJson(handlers.GetSuffixEngines()))
}

func TestPDFConverter(t *testing.T) {
	opts := []Option{
		AddMarkItDownConfig(&types.MarkItDownConfig{
			Command:    "markitdown",
			SafeDir:    "/Users/xxx/Downloads",
			TimeoutSec: 300,
		}),
		AddPandocConfig(&types.PandocConfig{
			Command:    "pandoc",
			SafeDir:    "/Users/xxx/Downloads",
			TimeoutSec: 300,
		}),
		AddTextInConfig(&types.TextInConfig{
			AppID:     "0839a4a63c591a19a8cca50643319287",
			AppSecret: "xxx",
			Host:      "https://api.textin.com",
		}),
		AddTongYiOcrConfig(&types.TongYiOcrConfig{
			ServerUrl: "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
			Secret:    "xxx",
			Model:     "qwen-vl-ocr",
		}),
	}

	c, err := New(log.DefaultLogger(), opts...)
	if err != nil {
		t.Fatal(err)
	}

	filepath := "./data/test.json"
	data, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}

	doc := &types.Document{
		Suffix:  filepath[strings.LastIndex(filepath, ".")+1:],
		Content: data,
	}
	result, err := c.Do(context.TODO(), doc, types.EnginePandoc, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(result))
}
