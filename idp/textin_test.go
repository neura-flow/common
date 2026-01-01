package idp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/neura-flow/common/util"
)

func TestPDF2MD(t *testing.T) {
	options := Options{
		PageStart:   util.Int(0),
		PageCount:   util.Int(1000), // 解析1000页
		TableFlavor: util.String("md"),
		ParseMode:   util.String("scan"), // 设置为scan模式
		Dpi:         util.Int(144),       // 分辨率为144 dpi
		PageDetails: util.Int(0),         // 不包含页面细节信息
	}
	kvm := make(map[string]interface{})
	_ = json.Unmarshal(util.ToJsonBytes(options), &kvm)
	fmt.Println(util.ToJson(kvm))

	textin := NewTextInOcr(
		&TextInConfig{
			AppID:     "0839a4a63c591a19a8cca50643319287",
			AppSecret: "xxx",
			Host:      "https://api.textin.com",
		},
	)

	image, err := getFileContent("./data/aws_text_test01.pdf")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	start := time.Now()
	fmt.Println("Request time: ", time.Now().Sub(start))

	doc := &Document{
		ContentType: "application/octet-stream",
		Content:     image,
		Options:     kvm,
	}
	result, err := textin.Convert(context.TODO(), doc)
	if err != nil {
		fmt.Println("Error convert:", err)
		return
	}
	//fmt.Printf("%s", string(result.Content))
	if err = writeFile(result.Content, fmt.Sprintf("./data/test%d.md", time.Now().UnixMilli())); err != nil {
		fmt.Println("Error writing file:", err)
	}
}

func writeFile(content []byte, filePath string) error {
	return os.WriteFile(filePath, content, 0644)
}
