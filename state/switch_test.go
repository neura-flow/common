package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwitch(t *testing.T) {
	h := SwitchHandler(func(s Switch, st State) {})
	s := NewSwitch(h)
	assert.True(t, s.On())
	assert.True(t, s.IsOn())
	assert.False(t, s.On())
	assert.True(t, s.Off())
	assert.False(t, s.Off())
	assert.False(t, s.IsOn())
}
