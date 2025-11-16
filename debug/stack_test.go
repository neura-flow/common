package debug

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCallerFrame(t *testing.T) {
	frame, ok := GetCallerFrame(0)
	assert.True(t, ok)
	ss := strings.Split(frame.Function, ".")
	assert.Equal(t, "TestGetCallerFrame", ss[len(ss)-1])
}

func TestGetStack(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			// 跳过 runtime 中的 defer 和 recover
			stack := GetStack(2, true)
			e := fmt.Sprintf("%v", err)
			prefix := "panic: "
			if strings.HasPrefix(e, "runtime") || strings.HasPrefix(e, "panic") {
				prefix = ""
			}
			fmt.Printf("%s%s\n\n%s", prefix, e, string(stack))
		}
	}()
	//testPanicA()
	testPanicB()
	panic("testA")
	*((*int)(nil)) = 1
}

func testPanicA() {
	testPanicB()
}

func testPanicB() {
	fmt.Println("testA")
	//panic("testA")
	*((*int)(nil)) = 1
}
