package debug

import (
	"bytes"
	"fmt"
	"runtime"
)

// GetCallerFrame gets caller frame. The argument skip is the number of stack
// frames to ascend, with 0 identifying the caller of getCallerFrame. The
// boolean ok is false if it was not possible to recover the information.
//
// Note: This implementation is similar to runtime.Caller, but it returns the whole frame.
func GetCallerFrame(skip int) (frame runtime.Frame, ok bool) {
	const skipOffset = 2 // skip getCallerFrame and Callers

	pc := make([]uintptr, 1)
	numFrames := runtime.Callers(skip+skipOffset, pc)
	if numFrames < 1 {
		return
	}

	frame, _ = runtime.CallersFrames(pc).Next()
	return frame, frame.PC != 0
}

// GetStack get stack and skip some callers
func GetStack(skip int, skipRuntime bool) []byte {
	stack := make([]byte, 4096)
	n := runtime.Stack(stack, true)
	if n <= 0 {
		return nil
	}
	stack = stack[:n]
	const skipOffset = 1 // skip GetStack
	frame, ok := GetCallerFrame(skip + skipOffset)
	if !ok {
		return stack
	}
	s := fmt.Sprintf("%s(", frame.Function)
	begin := bytes.Index(stack, []byte(s))
	if begin > 0 {
		stack = stack[begin:]
	}
	if skipRuntime {
		stack = SkipRuntime(stack)
	}
	return stack
}

// SkipRuntime skip runtime callers in stack
func SkipRuntime(stack []byte) []byte {
	// 找到以 panic( 开头的，意味着是 runtime 触发的 panic，需要跳过
	if begin := bytes.Index(stack, []byte("panic(")); (begin > 0 && (stack[begin-1] == '\n' || stack[begin-1] == '\r')) || begin == 0 {
		if (begin > 0 && (stack[begin-1] == '\n' || stack[begin-1] == '\r')) || begin == 0 {
			stack = stack[begin:]
		}
		if begin = bytes.Index(stack, []byte("src/runtime/panic.go:")); begin > 0 {
			for ; stack[begin] != '\r' && stack[begin] != '\n'; begin++ {
			}
			if stack[begin+1] == '\n' {
				begin++
			}
			stack = stack[begin:]
		}
	}

	return stack
}
