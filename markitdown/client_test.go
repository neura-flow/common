package markitdown

import (
	"fmt"
	"os"
	"testing"

	"github.com/neura-flow/common/log"
)

func TestDocxToMd(t *testing.T) {
	c, err := New(
		log.DefaultLogger(),
		&Config{
			Command: "markitdown",
			SafeDir: "/Users/xxx/Downloads",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("/Users/xx/Downloads/百炼系列手机产品介绍.docx")
	if err != nil {
		t.Fatal(err)
	}
	opts := Options{
		DataDir: "/Users/xxx/Downloads/markitdown",
	}
	doc, err := c.Convert(data, opts)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(doc))
}
