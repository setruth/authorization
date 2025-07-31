package util

import (
	"authorization.setruth.com/laa/exception"
	"authorization.setruth.com/laa/model"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/goccy/go-json"
	"os"
	"strings"
)

var saveFileName = "authCode"

func VerificationAuthCode(authCode string) (*model.AuthData, error) {
	//1. 分离密文和签名
	parts := strings.Split(authCode, ".")
	if len(parts) != 2 {
		return nil, exception.ErrInvalidAuthCodeFormat
	}

	encryptedAuthCodeBase64 := parts[0]
	signatureBase64 := parts[1]

	encryptedAuthCode, err := base64.StdEncoding.DecodeString(encryptedAuthCodeBase64)
	if err != nil {
		return nil, fmt.Errorf("%w: 密文base64解析失败: %v", exception.ErrInvalidAuthCodeFormat, err)
	}
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return nil, fmt.Errorf("%w: 签名base64解析失败: %v", exception.ErrInvalidAuthCodeFormat, err)
	}

	// 2. 验证数字签名
	hashed := sha256.Sum256(encryptedAuthCode)
	err = rsa.VerifyPKCS1v15(model.RSAPublicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return nil, exception.ErrSignatureVerificationFailed
	}

	// 3. 解密授权数据 (分离IV和密文)
	if len(encryptedAuthCode) < 12 {
		return nil, fmt.Errorf("%w: 密文IV长度不对", exception.ErrEncryptedIVParsingFailed)
	}
	iv := encryptedAuthCode[:12]
	encryptedAuthDetail := encryptedAuthCode[12:]

	block, err := aes.NewCipher(model.AESKey)
	if err != nil {
		return nil, fmt.Errorf("AES Cipher创建出错: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("AES GCM初始化失败: %v", err)
	}

	decryptedAuthDetailBytes, err := aesGCM.Open(nil, iv, encryptedAuthDetail, nil)
	if err != nil {
		return nil, exception.ErrEncryptedParsingFailed
	}
	// 3. 反序列化授权信息
	authDetail := &model.AuthData{}
	err = json.Unmarshal(decryptedAuthDetailBytes, authDetail)
	if err != nil {
		return nil, fmt.Errorf("授权码 JSON 解析失败: %v", err)
	}

	return authDetail, nil
}

func UpsertAuthCode(authCode string) error {
	contentBytes := []byte(authCode)
	err := os.WriteFile(saveFileName, contentBytes, 0644)
	if err != nil {
		return fmt.Errorf("授权码保存失败: %s", err.Error())
	}
	return nil
}
func ReadAuthCode() (string, error) {
	content, err := os.ReadFile(saveFileName)
	if err != nil {
		return "", fmt.Errorf("授权码读取失败: %s", err.Error())
	}
	return string(content), nil
}

func ClearAuthCode() {
	_ = os.WriteFile(saveFileName, []byte(""), 0644)
}
