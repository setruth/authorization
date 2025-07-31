package model

type AuthData struct {
	UniqueCode   string `json:"uniqueCode"`
	EndTimestamp int64  `json:"endTimestamp"`
}
type AddAuthRecordDTO struct {
	Name         string `json:"name"`
	UniqueCode   string `json:"uniqueCode"`
	EndTimestamp int64  `json:"endTimestamp"`
}
type UpdateEndTimestampDTO struct {
	ID           int   `json:"id"`
	EndTimestamp int64 `json:"endTimestamp"`
}

type Keys struct {
	RSAPrivateKey string `json:"rsaPrivateKey"`
	RSAPublicKey  string `json:"rsaPublicKey"`
	AESKey        string `json:"aesKey"`
}
