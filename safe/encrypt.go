package safe

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// 暂时只实现AES CBC 后续如果有需要新增的类型 考虑封装结构体

// 填充数据
func padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

// 去掉填充数据
func unpadding(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}

// EncryptAES
// ParamIN: src string 加密前的数据
// ParamIN: key string 加密密钥key
// ParamOUT: string 加密后的数据
// Notes: AES加密方法
func EncryptAES(src, key string) (string, error) {

	bytesKey := []byte(key)
	bytesSrc := []byte(src)
	block, err := aes.NewCipher(bytesKey)
	if err != nil {
		return "", err
	}

	bytesSrc = padding(bytesSrc, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, bytesKey)
	blockMode.CryptBlocks(bytesSrc, bytesSrc)
	return string(bytesSrc), nil
}

// DecryptAES
// ParamIN: dst string 加密后的数据
// ParamIN: key string 加密密钥key
// ParamOUT: string 解密后的数据
// Notes: AES解密方法
func DecryptAES(dst string, key string) (string, error) {

	bytesKey := []byte(key)
	bytesDst := []byte(dst)

	block, err := aes.NewCipher(bytesKey)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, bytesKey)
	blockMode.CryptBlocks(bytesDst, bytesDst)
	bytesDst = unpadding(bytesDst)
	return string(bytesDst), nil
}

// EncryptAESWithBase64
// 返回数据使用base64编码 返回乱码不方便数据库存储
func EncryptAESWithBase64(src, key string) (string, error) {

	encryptStr, err := EncryptAES(src, key)
	if err != nil {
		return "", err
	}

	sEnc := base64.StdEncoding.EncodeToString([]byte(encryptStr))
	return sEnc, nil
}

// DecryptAESWithBase64
// 同时解密需要先使用base64解码
func DecryptAESWithBase64(dst, key string) (string, error) {

	sDec, err := base64.StdEncoding.DecodeString(dst)
	if err != nil {
		return "", err
	}

	decryptStr, err := DecryptAES(string(sDec), key)
	if err != nil {
		return "", err
	}

	return decryptStr, nil
}
