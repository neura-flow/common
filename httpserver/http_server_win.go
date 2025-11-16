//go:build windows
// +build windows

package httpserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func initServer(address string, engine *gin.Engine) httpServer {
	s := &http.Server{
		Addr:           address,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return s
}
