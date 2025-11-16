package util

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
)

func GenSalt() (string, error) {
	bytes := make([]byte, 10)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", bytes), nil
}

func EncryptPwd(password, salt string) string {
	h := md5.New()
	io.WriteString(h, password+salt)
	return fmt.Sprintf("%x", h.Sum(nil))
}
