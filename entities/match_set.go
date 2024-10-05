package entities

import "github.com/google/uuid"

type MatchSet struct {
	Base
	MatchUUID  uuid.UUID `gorm:"type:CHAR(36)"`
	Number     int
	TeamAScore int
	TeamBScore int
}
