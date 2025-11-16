package mimetype

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/neura-flow/common/util"
)

func TestRead(t *testing.T) {
	data, err := os.ReadFile("./supported_mimes.md")
	if err != nil {
		t.Fatal(err)
	}

	kvm := map[string][]string{}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "| **") {
			continue
		}
		arr := strings.Split(line, "|")

		suffix := strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(arr[1]), "*", ""), ".", "")

		aliasesStr := strings.ReplaceAll(strings.TrimSpace(arr[3]), " ", "")
		var aliases []string
		if aliasesStr != "-" && aliasesStr != "" {
			aliases = strings.Split(aliasesStr, ",")
		} else {
			aliases = []string{}
		}

		mimeTypes := make([]string, 0)
		mimeType := strings.TrimSpace(arr[2])
		if mimeType != "" {
			mimeTypes = append(mimeTypes, mimeType)
		}
		mimeTypes = append(mimeTypes, aliases...)

		kvm[suffix] = mimeTypes
	}
	fmt.Println(util.ToJson(kvm))
}

func TestGetMimes(t *testing.T) {
	for k, v := range Get() {
		if len(v) > 1 {
			fmt.Printf("%s: %s\n", k, v)
		}
	}
}

func TestSuffixes(t *testing.T) {
	for k, v := range GetSuffixes() {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func TestDetectMimeType(t *testing.T) {
	content, err := os.ReadFile("/Users/liqj/Downloads/testdata/Go架构2024年度规划.pptx")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(Detect(content))
}

func TestDetectFile(t *testing.T) {
	fmt.Printf("%s\n", DetectFile("/Users/liqj/Downloads/testdata/Go架构2024年度规划.pptx"))
	fmt.Printf("%s\n", DetectFile("./supported_mimes.md"))
}

func TestDetectText(t *testing.T) {
	typ := Detect([]byte("<html>dd</html>"))
	fmt.Printf("%s\n", typ)
}
