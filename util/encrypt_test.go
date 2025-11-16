package util

import (
	"fmt"
	"testing"
)

func TestGenSalt(t *testing.T) {
	s, err := GenSalt()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", s)
}

func TestEncryptPwd(t *testing.T) {
	s := EncryptPwd("hello1234", "39f71b61012393be2962")
	fmt.Printf("%s\n", s)
}
