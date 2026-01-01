package markitdown

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
	Command    string `json:"command"`
	TimeoutSec int    `json:"timeoutSec"`
	SafeDir    string `json:"safeDir"`
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
		c.command = "markitdown"
	} else {
		c.command = cfg.Command
	}
	return c, nil
}

func (c *Client) Convert(data []byte, opts Options) ([]byte, error) {
	if len(opts.DataDir) > 0 && !strings.HasPrefix(opts.DataDir, c.safeDir) {
		return nil, fmt.Errorf("DataDir: '%s' is not in safe dir: '%s'", opts.DataDir, c.safeDir)
	}
	if err := c.checkDataDir(opts.DataDir); err != nil {
		return nil, err
	}
	tmpDir, err := os.MkdirTemp(opts.DataDir, "markitdown")
	if err != nil {
		return nil, err
	}
	defer remove(tmpDir)

	inputFile := filepath.Join(tmpDir, uuid.New())
	if err = os.WriteFile(inputFile, data, 0644); err != nil {
		return nil, err
	}
	defer remove(inputFile)

	var args []string
	var cleanups []func()
	if args, cleanups, err = opts.toCommandArgs(); err != nil {
		return nil, err
	}
	if len(cleanups) > 0 {
		defer func() {
			for i := 0; i < len(cleanups); i++ {
				cleanups[i]()
			}
		}()
	}

	var outputFile = filepath.Join(tmpDir, uuid.New())
	args = append(args, []string{inputFile, "--output", outputFile}...)
	if _, err = execCommand(c.logger, c.timeout, c.command, args...); err != nil {
		return nil, err
	}
	defer remove(outputFile)

	return os.ReadFile(outputFile)
}

func remove(file string) {
	_ = os.Remove(file)
}

func (c *Client) checkDataDir(dataDir string) error {
	if has, _ := pathExists(dataDir); !has {
		if err := createDir(c.logger, dataDir); err != nil {
			return err
		}
	}
	return nil
}

func pathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func createDir(logger log.Logger, dirs ...string) error {
	for _, v := range dirs {
		if exist, err := pathExists(v); err != nil {
			return err
		} else if !exist {
			if err = os.MkdirAll(v, os.ModePerm); err != nil {
				logger.Errorf("create directory: %s, err: %v", v, err)
				return err
			}
		}
	}
	return nil
}
