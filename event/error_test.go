package event

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := errors.New("testError")
	e := NewError(err, "test")
	assert.Equal(t, "test", string(e.Code()))
	assert.Equal(t, err, e.Unwrap())

	e = StringError("testError")
	assert.Equal(t, "testError", string(e.Code()))
	assert.Equal(t, "testError", e.Error())
}
