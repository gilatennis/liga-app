package mysql

import (
	"fmt"
	"liga-app/entities"

	"gorm.io/gorm"
)

func GormMigrate(db *gorm.DB) error {
	errMigrate := db.AutoMigrate(
		&entities.User{},
		&entities.MagicLinkLogin{},
		&entities.MagicLinkRegistration{},
		&entities.Wilayah{},
		&entities.League{},
		&entities.LeagueRegistration{},
		&entities.Team{},
		&entities.TeamPlayer{},
		&entities.LeaguePlayer{},
		&entities.LeagueTeam{},
		&entities.UserRating{},
		&entities.RuleBook{},
		&entities.MatchScheduleRequest{},
		&entities.Match{},
		&entities.MatchSet{},
		&entities.Venue{},
		&entities.Court{},
		&entities.CourtTimeSlot{},
		&entities.Group{},
		&entities.Fixture{},
		&entities.GroupTeam{},
		&entities.Round{},
		&entities.Admin{},
		&entities.GroupRoundRobinTeamResult{},
		&entities.GroupRoundRobinHeadToHeadResult{},
	)

	if errMigrate != nil {
		return fmt.Errorf("[mysql.Migrate]: error migrating: %w", errMigrate)
	}

	return nil
}
