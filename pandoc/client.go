package pandoc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/neura-flow/common/log"

	"github.com/pborman/uuid"
)

type Config struct {
	Command         string `json:"command"`
	TimeoutSec      int    `json:"timeoutSec"`
	Verbose         bool   `json:"verbose"`
	Trace           bool   `json:"trace"`
	DumpArgs        bool   `json:"dumpArgs"`
	IgnoreArgs      bool   `json:"ignoreArgs"`
	EnableFilter    bool   `json:"enableFilter"`
	EnableLuaFilter bool   `json:"enableLuaFilter"`
	SafeDir         string `json:"safeDir"`
}

type Client struct {
	logger log.Logger
	cfg    *Config

	command string
	timeout time.Duration
	safeDir string
}

func New(logger log.Logger, cfg *Config) (*Client, error) {
	c := &Client{
		logger: logger,
		cfg:    cfg,
	}
	if cfg.TimeoutSec > 0 {
		c.timeout = time.Second * time.Duration(cfg.TimeoutSec)
	} else {
		c.timeout = time.Second * 300
	}
	if cfg.SafeDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		c.safeDir = cwd
	} else {
		c.safeDir = cfg.SafeDir
	}
	if cfg.Command == "" {
		c.command = "pandoc"
	} else {
		c.command = cfg.Command
	}
	return c, nil
}

func (c *Client) Convert(data []byte, opts Options) (ret []byte, err error) {
	if len(opts.DataDir) > 0 && !strings.HasPrefix(opts.DataDir, c.safeDir) {
		err = fmt.Errorf("DataDir: '%s' is not in safe dir: '%s'", opts.DataDir, c.safeDir)
		return
	}

	tmpDir, err := os.MkdirTemp(opts.DataDir, "go-pandoc")
	if err != nil {
		return
	}
	defer remove(tmpDir)

	inputFile := filepath.Join(tmpDir, uuid.New()) + "." + opts.From
	if err = os.WriteFile(inputFile, data, 0644); err != nil {
		return
	}
	defer remove(inputFile)

	opts.verbose = c.cfg.Verbose
	opts.trace = c.cfg.Trace
	opts.dumpArgs = c.cfg.DumpArgs
	opts.ignoreArgs = c.cfg.IgnoreArgs

	var args []string
	var cleanups []func()
	if args, cleanups, err = opts.toCommandArgs(c.safeDir, c.cfg.EnableFilter, c.cfg.EnableLuaFilter); err != nil {
		return
	}
	if len(cleanups) > 0 {
		defer func() {
			for i := 0; i < len(cleanups); i++ {
				cleanups[i]()
			}
		}()
	}

	var outputFile = filepath.Join(tmpDir, uuid.New()) + "." + opts.To
	args = append(args, []string{"--quiet", inputFile, "--output", outputFile}...)
	if _, err = execCommand(c.logger, c.timeout, c.command, args...); err != nil {
		return
	}
	defer remove(outputFile)

	return os.ReadFile(outputFile)
}

func toArgs(k, url, safeDir string) ([]string, func(), error) {
	var file = File{
		Url:           url,
		TempDirPrefix: "go-pandoc",
		SafeDir:       safeDir,
	}
	var filename string
	var err error
	if filename, err = file.Path(); err != nil {
		return nil, nil, err
	}

	args := []string{k, filename}
	return args, file.Cleanup, nil
}

func remove(file string) {
	_ = os.Remove(file)
}
