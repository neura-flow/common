package tongyi_file

import (
	"fmt"
	"os"
	"testing"

	"github.com/neura-flow/common/util"
)

func TestCreate(t *testing.T) {
	client, err := New(&Config{
		Secret: "xxx",
	})
	if err != nil {
		panic(err)
	}

	filename := "百炼系列手机产品介绍.docx"
	data, err := os.ReadFile("./data/百炼系列手机产品介绍.docx")
	if err != nil {
		panic(err)
	}

	params := map[string]string{
		"purpose": "file-extract",
	}
	resp, err := client.Create(filename, params, data)
	if err != nil {
		panic(err)
	}

	// {
	//	"file": {
	//		"id": "file-fe-1pWjv6qZ4ShjWzEMOJzU9DoD",
	//		"bytes": 0,
	//		"created_at": 1741574115,
	//		"filename": "百炼系列手机产品介绍.docx",
	//		"object": "file",
	//		"purpose": "file-extract",
	//		"status": "processed"
	//	},
	//	"error": null
	//}
	fmt.Printf("%s\n", util.ToJson(resp))
}
