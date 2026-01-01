package gosseract

import (
	"fmt"
	"testing"

	"github.com/otiai10/gosseract/v2"
)

func TestExec(t *testing.T) {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage("./data/001-helloworld.png")
	text, _ := client.Text()
	fmt.Println(text)
}
