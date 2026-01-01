package tongyi

import (
	"context"
	"fmt"
	"testing"

	"github.com/neura-flow/common/embedding/types"
	"github.com/neura-flow/common/util"
)

func TestCall(t *testing.T) {
	serverUrl := "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding"
	apiKey := "xxx"

	client := New(&Config{
		ServerUrl: serverUrl,
		ApiKey:    apiKey,
	})

	embRequest := &types.Request{
		Model: "text-embedding-v3",
		Input: &types.RequestInput{
			Texts: []string{"衣服的质量杠杠的，很漂亮，不枉我等了这么久啊，喜欢，以后还来这里买"},
		},
		Parameters: &types.RequestParameters{
			Dimension: 1024,
		},
	}
	fmt.Printf("request: %s\n", util.ToJson(embRequest))

	embResp, err := client.Call(context.TODO(), embRequest)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", util.ToJson(embResp))
}
