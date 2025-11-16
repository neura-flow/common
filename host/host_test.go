package host

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHost(t *testing.T) {
	addr := "127.0.0.1:12345"
	h, p, err := ExtractHostPort(addr)
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", h)
	assert.Equal(t, uint64(12345), p)
	l, err := net.Listen("tcp", addr)
	assert.NoError(t, err)
	assert.NotNil(t, l)
	port, ok := Port(l)
	assert.True(t, ok)
	addr1, err := Extract(addr, nil)
	assert.NoError(t, err)
	assert.Equal(t, addr, addr1)
	assert.Equal(t, 12345, port)
	addr1, err = Extract(addr, l)
	assert.NoError(t, err)
	assert.Equal(t, addr, addr1)
	addr1, err = Extract("", l)
	assert.NoError(t, err)
	h, p, err = ExtractHostPort(addr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(12345), p)
	l.Close()
	ip, err := IP()
	assert.NoError(t, err)
	assert.NotEqual(t, ip, "")
}
