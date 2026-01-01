package safe

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	s := String{"test"}
	data, err := json.Marshal(&s)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(data))
	var s1 String
	err = json.Unmarshal([]byte("\"test\""), &s1)
	assert.NoError(t, err)
	assert.Equal(t, "test", s1.Value())
	assert.Equal(t, "", s1.String())
	s1.SetValue("abc")
	assert.Equal(t, "abc", s1.Value())
	assert.Equal(t, "", s1.String())

	type testData struct {
		Test String `json:"test"`
	}
	data = []byte(`{"test":"hello"}`)
	var d testData
	err = json.Unmarshal(data, &d)
	assert.NoError(t, err)
	data, err = json.Marshal(d)
	assert.NoError(t, err)
	t.Log(string(data))
	ss := fmt.Sprint(d.Test)
	assert.Equal(t, "", ss)
}
