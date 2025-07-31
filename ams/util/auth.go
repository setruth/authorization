package util

import (
	"authorization.setruth.com/ams/model"
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/goccy/go-json"
)

func GenerateAuthorizationCode(uniqueCode string, endTimestamp int64) (string, error) {
	rsaPrivateKey, err := model.GetRsaPrivateKey()
	if err != nil {
		return "", err
	}

	// 1. 构造原始授权数据并序列化为 JSON
	authData := model.AuthData{
		UniqueCode:   uniqueCode,
		EndTimestamp: endTimestamp,
	}
	authDetailBytes, err := json.Marshal(authData)
	if err != nil {
		return "", fmt.Errorf("序列化授权数据到 JSON 失败: %w", err)
	}

	// 2. 对称加密 AuthData (AES-GCM)
	iv := make([]byte, 12)
	if _, err := rand.Read(iv); err != nil {
		return "", fmt.Errorf("生成 AES-GCM IV 失败: %w", err)
	}
	block, err := aes.NewCipher(model.AESKey)
	if err != nil {
		return "", fmt.Errorf("AES Cipher创建出错: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建 GCM 加密器失败: %w", err)
	}
	encryptedAuthDetailWithTag := gcm.Seal(nil, iv, authDetailBytes, nil)

	encryptedAuthDetailCombined := bytes.Join([][]byte{iv, encryptedAuthDetailWithTag}, nil)

	// 3. 对组合后的加密数据进行 SHA256 哈希并使用 RSA 私钥签名
	hashed := sha256.Sum256(encryptedAuthDetailCombined)
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashed[:]) // hashed[:] 将 [32]byte 转换为 []byte
	if err != nil {
		return "", fmt.Errorf("对加密数据签名失败: %w", err)
	}
	// 4. 组装授权码：Base64 编码的加密数据 + "." + Base64 编码的签名
	encodedEncryptedData := base64.StdEncoding.EncodeToString(encryptedAuthDetailCombined)
	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	return fmt.Sprintf("%s.%s", encodedEncryptedData, encodedSignature), nil
}
