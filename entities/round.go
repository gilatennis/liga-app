package entities

import (
	"time"

	"github.com/google/uuid"
)

type Round struct {
	Base
	GroupUUID uuid.UUID `gorm:"type:CHAR(36)"`
	Number    int
	Matches   []Match `gorm:"ForeignKey:"RoundUUID"`
	Start     *time.Time
	End       *time.Time
}
