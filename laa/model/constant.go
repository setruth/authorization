package model

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log"
	"sync"
)

var (
	rsaPublicKeyBASE64 = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAk1/P/VaorDFGw7DdVHNuCIek7TieJEu7kchbER5fZNan0+lPg5TXQNef8KgtMLbGb63odz0EYrutXLcZjGqEJwBNCKpzyHG2EObrKjBZjhoG+KmbRou1hIEX/FCknaHZm1B5uDqHDUrcHpvXZYftXGXfCv5lCs9s4VWcNFMm/5C9OkPb1zku++2unOq/9qaUMv8M4FWMLEWOnYCkJqrjR57ll1Ys7BGhDKqM7A4oQ1l/04nt96Y8q6A4+2UGOntuE/Rfu5wJKnF4KhRf+36smfTUbNAR31ftM7vMVoLCbqGpLXnlf9+7J8obEJZn0paUH2KACeHTPGlM3u4P3jXMVwIDAQAB"
	aesKeyBASE64       = "lBbYXsW3hbBc6IyUOXAPelaWB7t+lsqLbzyaO1oM+uU="
	AESKey             []byte
	RSAPublicKey       *rsa.PublicKey

	TaskStopChan chan struct{}
	TaskWg       sync.WaitGroup

	UniqueCodeCache string
	AuthDetailCache *AuthDetail = nil
)

const EndTimestampNil = -1

type AuthStatusTag int

const (
	// Unauthorized 未授权
	Unauthorized AuthStatusTag = iota
	// Authorized 授权成功
	Authorized
	// Expire 授权到期
	Expire
)

func (s AuthStatusTag) String() string {
	switch s {
	case Unauthorized:
		return "UNAUTHORIZED"
	case Authorized:
		return "AUTHORIZED"
	case Expire:
		return "EXPIRE"
	default:
		return fmt.Sprintf("AuthStatusTag(%d)", s)
	}
}
func init() {
	TaskStopChan = make(chan struct{})
	// 初始化 AESKey
	decodedAESKey, err := base64.StdEncoding.DecodeString(aesKeyBASE64)
	if err != nil {
		log.Fatalf("你的AES密钥不是标准的BASE64编码内容无法解码: %v", err)
	}
	AESKey = decodedAESKey
	if len(AESKey) != 16 && len(AESKey) != 24 && len(AESKey) != 32 {
		log.Fatalf("ASE密钥的长度:%d有问题,不是标准的AES密钥长度", len(AESKey))
	}

	// 初始化 RSAPublicKey
	decodedRSAPublicKeyBytes, err := base64.StdEncoding.DecodeString(rsaPublicKeyBASE64)
	if err != nil {
		log.Fatalf("你的AES密钥不是标准的BASE64编码内容无法解码: %v", err)
	}

	// 尝试解析为 PKIX 格式的公钥（X.509 SubjectPublicKeyInfo）
	pub, err := x509.ParsePKIXPublicKey(decodedRSAPublicKeyBytes)
	if err != nil {

		log.Fatalf("解析公钥失败:%v", len(AESKey))
	}

	var ok bool
	RSAPublicKey, ok = pub.(*rsa.PublicKey)
	if !ok {
		log.Fatal("进行断言发现并不是PublicKey类型")
	}
}
