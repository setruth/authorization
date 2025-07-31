package model

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"log"
	"sync"
)

var (
	aesKeyBASE64       = "lBbYXsW3hbBc6IyUOXAPelaWB7t+lsqLbzyaO1oM+uU="
	AESKey             []byte
	rsaPrivateKey      *rsa.PrivateKey = nil
	rsaPrivateKeyMutex sync.RWMutex
)

const RsaPrivateEnvKey = "AUTHORIZATION_SYSTEM_RSA_PRIVATE_KEY"

func UpdateRsaPrivateKey(key *rsa.PrivateKey) {
	rsaPrivateKeyMutex.Lock()
	defer rsaPrivateKeyMutex.Unlock()
	rsaPrivateKey = key
	log.Println("RSA 私钥已更新。")
}
func GetRsaPrivateKey() (*rsa.PrivateKey, error) {
	rsaPrivateKeyMutex.RLock()         // 获取读锁
	defer rsaPrivateKeyMutex.RUnlock() // 函数退出时释放读锁

	if rsaPrivateKey == nil {
		return nil, errors.New("RSA 私钥未加载")
	}
	return rsaPrivateKey, nil
}
func init() {
	// 初始化 AESKey
	decodedAESKey, err := base64.StdEncoding.DecodeString(aesKeyBASE64)
	if err != nil {
		log.Fatalf("你的AES密钥不是标准的BASE64编码内容无法解码: %v", err)
	}
	AESKey = decodedAESKey
	if len(AESKey) != 16 && len(AESKey) != 24 && len(AESKey) != 32 {
		log.Fatalf("ASE密钥的长度:%d有问题,不是标准的AES密钥长度", len(AESKey))
	}
}
