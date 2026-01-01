package markitdown

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"syscall"
	"time"

	"github.com/neura-flow/common/log"
)

func execCommand(logger log.Logger, timeout time.Duration, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	var err error
	var stdin io.WriteCloser
	var stdout, stderr io.ReadCloser
	if stdin, err = cmd.StdinPipe(); err != nil {
		return nil, err
	}
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return nil, err
	}
	if stderr, err = cmd.StderrPipe(); err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, err
	}
	_ = stdin.Close()

	outBuf := bytes.NewBuffer(nil)
	errBuf := bytes.NewBuffer(nil)
	go func() {
		if _, err = io.Copy(errBuf, stderr); err != nil {
			logger.Errorf("Failed to copy stderr to stdout: %s", err)
		}
	}()
	go func() {
		if _, err = io.Copy(outBuf, stdout); err != nil {
			logger.Errorf("Failed to copy stdout to stdout: %s", err)
		}
	}()

	ch := make(chan error)
	go func(cmd *exec.Cmd) {
		defer close(ch)
		ch <- cmd.Wait()
	}(cmd)

	select {
	case err0, ok := <-ch:
		if ok && err0 != nil {
			logger.Errorf("%v", err0)
			return nil, errors.New(errBuf.String())
		}
	case <-time.After(timeout):
		if err = cmd.Process.Kill(); err != nil {
			logger.Errorf("Failed to kill process: %s", err)
		}
		return nil, errors.New("execute timeout")
	}

	return outBuf.Bytes(), nil
}
