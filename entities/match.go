package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Match struct {
	Base
	Number                   int
	LeagueUUID               uuid.UUID  `gorm:"type:CHAR(36)"`
	League                   League     `gorm:"ForeignKey:LeagueUUID;References:UUID"`
	RoundUUID                uuid.UUID  `gorm:"type:CHAR(36)"`
	Round                    Round      `gorm:"ForeignKey:RoundUUID;References:UUID"`
	TeamAUUID                *uuid.UUID `gorm:"type:CHAR(36)"`
	TeamA                    *Team      `gorm:"ForeignKey:TeamAUUID;References:UUID"`
	TeamBUUID                *uuid.UUID `gorm:"type:CHAR(36)"`
	TeamB                    *Team      `gorm:"ForeignKey:TeamBUUID;References:UUID"`
	IsByeTeamA               bool
	IsByeTeamB               bool
	Location                 string
	TimeStart                *time.Time
	TimeEnd                  *time.Time
	Sets                     []MatchSet `gorm:"ForeignKey:MatchUUID"`
	Status                   string
	FilterStatus             string                `gorm:"-:all"`
	MatchScheduleRequestUUID *uuid.UUID            `gorm:"type:CHAR(36)"`
	MatchScheduleRequest     *MatchScheduleRequest `gorm:"ForeignKey:MatchScheduleRequestUUID;References:UUID;constraint:OnDelete:SET NULL;"`
	IsComplete               bool
	IsCompletedBy            uuid.UUID `gorm:"type:CHAR(36)"`
	IsCompletedAt            *time.Time
	IsCompleteConfirmed      bool
	IsCompleteConfirmedBy    uuid.UUID `gorm:"type:CHAR(36)"`
	IsCompleteConfirmedAt    *time.Time
}

func (m *Match) GetMatchSetByNumber(num int) *MatchSet {
	var ms *MatchSet
	for _, s := range m.Sets {
		if s.Number == num {
			ss := s
			ms = &ss
		}
	}

	return ms
}

func (m *Match) GetMatchSetsScoreString() string {
	scores := []string{}
	for _, s := range m.Sets {
		score := fmt.Sprintf("%d - %d", s.TeamAScore, s.TeamBScore)
		scores = append(scores, score)
	}

	return strings.Join(scores, ", ")
}
