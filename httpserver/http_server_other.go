//go:build !windows
// +build !windows

package httpserver

import (
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func initServer(address string, engine *gin.Engine) httpServer {
	s := endless.NewServer(address, engine)
	s.ReadHeaderTimeout = 15 * time.Second
	s.WriteTimeout = 180 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}
