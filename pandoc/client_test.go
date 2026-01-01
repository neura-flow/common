package pandoc

import (
	"fmt"
	"os"
	"testing"

	"github.com/neura-flow/common/log"
)

func TestCsvToMd(t *testing.T) {
	c, err := New(
		log.DefaultLogger(),
		&Config{
			Command: "pandoc",
			SafeDir: "./data",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("./data/2016-07-29.csv")
	if err != nil {
		t.Fatal(err)
	}
	opts := Options{
		From:    "csv",
		To:      "markdown_mmd",
		DataDir: "./data/pandoc_data",
	}
	doc, err := c.Convert(data, opts)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(doc))
}
