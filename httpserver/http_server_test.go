package httpserver_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/neura-flow/common/httpserver"
	"github.com/neura-flow/common/log"
)

func TestHttpServer(t *testing.T) {
	logger := log.DefaultLogger()
	svr := httpserver.NewHttpServer(logger, &httpserver.Config{
		GinMode:    gin.DebugMode,
		ServerPort: 10001,
	})
	if err := svr.Start(); err != nil {
		t.Fatal(err)
	}
}
