package cipher

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
)

const (
	RSABits1024    = 1024
	RSABits2048    = 2048
	TypePublicKey  = "PUBLIC KEY"
	TypePrivateKey = "RSA PRIVATE KEY"
	FormatPkcs1    = "pkcs1"
	FormatPkcs8    = "pkcs8"
)

type RSAConfig struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
	// 支持格式: pkcs1,pkcs8, 默认: pkcs1
	Format string `json:"format"`
}

type RSA struct {
	cfg *RSAConfig
}

func NewRSA(cfg *RSAConfig) *RSA {
	return &RSA{
		cfg: cfg,
	}
}

func (r *RSA) pkcs8() bool {
	return strings.EqualFold("pkcs8", r.cfg.Format)
}

func (r *RSA) EncryptToBase64(plain string) (string, error) {
	bytes, err := r.Encrypt([]byte(plain))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func (r *RSA) DecryptBase64(cipher string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(cipher)
	if err != nil {
		return "", err
	}
	ret, err := r.Decrypt(bytes)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

func (r *RSA) Encrypt(plain []byte) ([]byte, error) {
	block, rest := pem.Decode([]byte(r.cfg.PublicKey))
	if block == nil || block.Type != TypePublicKey {
		return nil, fmt.Errorf("failed to decode pem block containing public key, rest: %s", string(rest))
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pubKey := pub.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, plain)
}

func (r *RSA) Decrypt(cipher []byte) ([]byte, error) {
	block, rest := pem.Decode([]byte(r.cfg.PrivateKey))
	if block == nil || block.Type != TypePrivateKey {
		return nil, fmt.Errorf("failed to decode pem block containing private key, rest: %s", string(rest))
	}
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, pri, cipher)
}

func GenRSA(bits int, pkcs8 bool) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}

	var priBlock []byte
	if pkcs8 {
		priBlock, err = x509.MarshalPKCS8PrivateKey(privateKey)
	} else {
		priBlock = x509.MarshalPKCS1PrivateKey(privateKey)
	}
	if err != nil {
		return nil, nil, err
	}

	priBytes := pem.EncodeToMemory(&pem.Block{
		Type:  TypePrivateKey,
		Bytes: priBlock,
	})
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  TypePublicKey,
		Bytes: derPkix,
	})
	return pubBytes, priBytes, nil
}
