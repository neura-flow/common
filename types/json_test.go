package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJson(t *testing.T) {
	var v string
	var err error
	var b Bool
	v = "\"true\""
	err = json.Unmarshal([]byte(v), &b)
	assert.NoError(t, err)
	assert.Equal(t, true, b.Val())
	v = "false"
	err = json.Unmarshal([]byte(v), &b)
	assert.NoError(t, err)
	assert.Equal(t, false, b.Val())
	var i Int
	v = "12345"
	err = json.Unmarshal([]byte(v), &i)
	assert.NoError(t, err)
	assert.Equal(t, 12345, i.Val())
	v = "\"12345\""
	err = json.Unmarshal([]byte(v), &i)
	assert.NoError(t, err)
	assert.Equal(t, 12345, i.Val())
	var f Float
	v = "12345"
	err = json.Unmarshal([]byte(v), &f)
	assert.NoError(t, err)
	assert.Equal(t, float64(12345), f.Val())
	v = "\"1234.5\""
	err = json.Unmarshal([]byte(v), &f)
	assert.NoError(t, err)
	assert.Equal(t, float64(1234.5), f.Val())
}
