package host

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemote(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	addr := "127.0.0.1"
	req.RemoteAddr = addr + ":12345"
	assert.Equal(t, addr, RemoteIP(req))
	addr = "127.0.0.2"
	req.Header.Set(XForwardedFor, addr)
	assert.Equal(t, addr, RemoteIP(req))
	addr = "127.0.0.3"
	req.Header.Set(XRealIP, addr)
	assert.Equal(t, addr, RemoteIP(req))
}
