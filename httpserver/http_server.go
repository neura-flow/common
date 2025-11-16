package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neura-flow/common/log"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

type Config struct {
	GinMode    string `json:"ginMode,omitempty"`
	ServerPort int    `json:"serverPort,omitempty"`
}

type httpServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type HttpServer struct {
	logger  log.Logger
	cfg     *Config
	server  httpServer
	engine  *gin.Engine
	running int32
	stopCh  chan interface{}
}

func NewHttpServer(logger log.Logger, cfg *Config) *HttpServer {
	if cfg.GinMode == "" {
		cfg.GinMode = gin.ReleaseMode
	}
	s := &HttpServer{
		logger: logger,
		cfg:    cfg,
		stopCh: make(chan interface{}),
	}
	s.setGinMode()
	s.engine = s.createEngine()
	return s
}

func (s *HttpServer) HandlePrefix(prefix string, fn http.HandlerFunc) {
	s.engine.Any(fmt.Sprintf("%s/*name", prefix), func(c *gin.Context) {
		fn(c.Writer, c.Request)
	})
}

func (s *HttpServer) Start() error {
	if !atomic.CompareAndSwapInt32(&s.running, 0, 1) {
		return nil
	}

	s.server = initServer(fmt.Sprintf(":%d", s.cfg.ServerPort), s.engine)

	s.logger.Infof("starting http server, listening on: http://127.0.0.1:%d", s.cfg.ServerPort)

	if err := s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) ||
			strings.Contains(err.Error(), "use of closed network connection") {
			return nil
		} else {
			s.logger.Errorf("%v", err)
			return err
		}
	}
	return nil
}

func (s *HttpServer) Stop() error {
	if !atomic.CompareAndSwapInt32(&s.running, 1, 0) {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *HttpServer) createEngine() *gin.Engine {
	var engine = gin.Default()
	ginMonitor().Use(engine)
	engine.ForwardedByClientIP = true
	engine.Use(ginRecovery(s.logger, true))
	engine.RemoteIPHeaders = []string{"X-Forwarded-For", "X-Real-IP"}
	engine.GET("health", s.health)
	return engine
}

func (s *HttpServer) setGinMode() {
	switch s.cfg.GinMode {
	case gin.ReleaseMode:
		gin.SetMode(gin.ReleaseMode)
	case gin.TestMode:
		gin.SetMode(gin.TestMode)
	case gin.DebugMode:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
}

func (s *HttpServer) health(c *gin.Context) {
	c.Header("Connection", "close")
	c.Request.Close = true
	c.JSON(http.StatusOK, map[string]string{
		"Status":      "UP",
		"Description": "",
	})
}

func (s *HttpServer) Engine() *gin.Engine {
	return s.engine
}

func ginMonitor() *ginmetrics.Monitor {
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(3)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	return m
}

func ginRecovery(logger log.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Errorf("url: %s request: %s err: %v ", c.Request.URL.Path, string(httpRequest), err)
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Errorf("[Recovery from panic] request: %s stack: %s err: %v", string(httpRequest), string(debug.Stack()), err)
				} else {
					logger.Errorf("[Recovery from panic] request: %s err: %v", string(httpRequest), err)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
