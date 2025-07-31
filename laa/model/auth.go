package model

type AuthData struct {
	UniqueCode   string `json:"uniqueCode"`
	EndTimestamp int64  `json:"endTimestamp"`
}
type AuthStatus struct {
	Tag          AuthStatusTag `json:"tag"`
	EndTimestamp int64         `json:"endTimestamp"`
}

type AuthDetail struct {
	AuthCode     string
	UniqueCode   string
	EndTimestamp int64
}
