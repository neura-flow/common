package tongyi_ocr

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/neura-flow/common/llm/types"
	"github.com/neura-flow/common/util"
)

func TestPwd(t *testing.T) {
	path, _ := os.Getwd()
	fmt.Printf("%s\n", path)
}

func TestCall(t *testing.T) {
	client := New(&Config{
		ServerUrl: "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
		//ServerUrl: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		Secret: "xxx",
	})

	//1. read image url
	imageUrl := "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/ctdzex/biaozhun.jpg"
	completion := types.NewTongYiCompletionRequest("qwen-vl-ocr")
	completion.Messages = []any{
		types.UserMessage{
			Role: "user", Content: NewImageUrlMessageContents(imageUrl),
		},
	}

	// 2. read image file
	//path, _ := os.Getwd()
	//data, err := os.ReadFile(fmt.Sprintf("%s/data/biaozhun.jpg", path))
	//if err != nil {
	//	panic(err)
	//}
	//completion := types.NewTongYiCompletionRequest("qwen-vl-ocr")
	//completion.Messages = []any{
	//	types.UserMessage{
	//		Role: "user", Content: NewImageBlobMessageContents("image/jpeg", data),
	//	},
	//}

	resp, err := client.Call(context.TODO(), completion)
	if err != nil {
		panic(err)
	}
	fmt.Println("Response:")
	fmt.Printf("%s\n", util.ToJson(resp))
}
