package exception

import "errors"

var (
	ErrInvalidAuthCodeFormat       = errors.New("不合法的授权码格式")
	ErrSignatureVerificationFailed = errors.New("授权码数字签名验证失败")
	ErrEncryptedParsingFailed      = errors.New("授权码解析失败")
	ErrEncryptedIVParsingFailed    = errors.New("授权密文IV分离失败")
)
