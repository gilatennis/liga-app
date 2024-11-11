package entities

import "github.com/google/uuid"

type MatchStats struct {
	UserUUID uuid.UUID `gorm:"-:migration"`
	User     User
	Played   int `gorm:"-:migration"`
	Win      int `gorm:"-:migration"`
	Lose     int `gorm:"-:migration"`
}
