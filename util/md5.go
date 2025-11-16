package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

func MD5(input string) string {
	c := md5.New()
	c.Write([]byte(input))
	bytes := c.Sum(nil)
	return hex.EncodeToString(bytes)
}

// EncodeMD5 对字符串进行md5处理
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

// EncodeSha1 对字符串进行sha1处理
func EncodeSha1(value string) string {
	h := sha1.New()
	h.Write([]byte(value))

	return hex.EncodeToString(h.Sum(nil))
}

func HmacSHA256(plaintext string, key string) string {
	hash := hmac.New(sha256.New, []byte(key)) // 创建哈希算法
	hash.Write([]byte(plaintext))             // 写入数据
	return fmt.Sprintf("%X", hash.Sum(nil))
}

func HmacMD5(plaintext string, key string) string {
	hash := hmac.New(md5.New, []byte(key)) // 创建哈希算法
	hash.Write([]byte(plaintext))          // 写入数据
	return fmt.Sprintf("%X", hash.Sum(nil))
}
