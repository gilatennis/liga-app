package entities

import "github.com/google/uuid"

type UserRating struct {
	Base
	UserUUID uuid.UUID `gorm:"type:CHAR(36);uniqueIndex:idx_user_rating_season"`
	// User       User      `gorm:"ForeignKey:UserUUID;References:UUID"`
	NTRPRating float32
	Season     int `gorm:"uniqueIndex:idx_user_rating_season"`
}
