package entities

import (
	"liga-app/enums"

	"github.com/google/uuid"
)

type Fixture struct {
	Base
	LeagueUUID  uuid.UUID
	League      League `gorm:"ForeignKey:LeagueUUID;References:UUID"`
	Name        string
	StageNumber int
	Format      enums.FIXTURE_FORMAT
	Groups      []Group `gorm:"ForeignKey:FixtureUUID"`
}
