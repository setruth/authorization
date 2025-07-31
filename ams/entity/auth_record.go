package entity

type AuthRecord struct {
	ID           int    `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	UniqueCode   string `gorm:"type:varchar(255);uniqueIndex;not null" json:"uniqueCode"`
	EndTimestamp int64  `gorm:"not null" json:"endTimestamp"`
}
