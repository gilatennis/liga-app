package leagues

import (
	"errors"
	"fmt"
	"liga-app/entities"
	"liga-app/enums"
	"liga-app/utils"

	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

var DEFAULT_LEAGUES []entities.League = []entities.League{
	{
		Name:                  "Beregu Putri",
		WilayahIDProvince:     "31",
		ParticipationCategory: enums.CATEGORY_PARTICIPATION_TEAMS,
		GenderCategory:        enums.CATEGORY_GENDER_FEMALE,
		Level:                 enums.CATEGORY_LEVEL_BEGINNER,
		LinkWAGroup:           "https://chat.whatsapp.com/JUup3ZHCEMd1eO2WifCsI5",
		Season:                1,
		RegistrationFee:       1_500_000,
		LinkWARatingAdmin:     "628159239794",
	},
	{
		Name:                  "Beregu Putra",
		WilayahIDProvince:     "31",
		ParticipationCategory: enums.CATEGORY_PARTICIPATION_TEAMS,
		GenderCategory:        enums.CATEGORY_GENDER_MALE,
		Level:                 enums.CATEGORY_LEVEL_BEGINNER,
		LinkWAGroup:           "https://chat.whatsapp.com/JUup3ZHCEMd1eO2WifCsI5",
		Season:                1,
		RegistrationFee:       1_500_000,
		LinkWARatingAdmin:     "628159239794",
	},
	{
		Name:                  "Singles Putra",
		WilayahIDProvince:     "31",
		ParticipationCategory: enums.CATEGORY_PARTICIPATION_SINGLES,
		GenderCategory:        enums.CATEGORY_GENDER_MALE,
		Level:                 enums.CATEGORY_LEVEL_BEGINNER,
		LinkWAGroup:           "https://chat.whatsapp.com/JUup3ZHCEMd1eO2WifCsI5",
		Season:                1,
		RegistrationFee:       200_000,
		LinkWARatingAdmin:     "628159239794",
	},
	{
		Name:                  "Singles Putri",
		WilayahIDProvince:     "31",
		ParticipationCategory: enums.CATEGORY_PARTICIPATION_SINGLES,
		GenderCategory:        enums.CATEGORY_GENDER_FEMALE,
		Level:                 enums.CATEGORY_LEVEL_BEGINNER,
		LinkWAGroup:           "https://chat.whatsapp.com/JUup3ZHCEMd1eO2WifCsI5",
		Season:                1,
		RegistrationFee:       200_000,
		LinkWARatingAdmin:     "628159239794",
	},
	{
		Name:                  "Doubles Putra",
		WilayahIDProvince:     "31",
		ParticipationCategory: enums.CATEGORY_PARTICIPATION_DOUBLES,
		GenderCategory:        enums.CATEGORY_GENDER_MALE,
		Level:                 enums.CATEGORY_LEVEL_BEGINNER,
		LinkWAGroup:           "https://chat.whatsapp.com/JUup3ZHCEMd1eO2WifCsI5",
		Season:                1,
		RegistrationFee:       400_000,
		LinkWARatingAdmin:     "628159239794",
	},
	{
		Name:                  "Doubles Putri",
		WilayahIDProvince:     "31",
		ParticipationCategory: enums.CATEGORY_PARTICIPATION_DOUBLES,
		GenderCategory:        enums.CATEGORY_GENDER_FEMALE,
		Level:                 enums.CATEGORY_LEVEL_BEGINNER,
		LinkWAGroup:           "https://chat.whatsapp.com/JUup3ZHCEMd1eO2WifCsI5",
		Season:                1,
		RegistrationFee:       400_000,
		LinkWARatingAdmin:     "628159239794",
	},
	{
		Name:                  "Doubles Mixed",
		WilayahIDProvince:     "31",
		ParticipationCategory: enums.CATEGORY_PARTICIPATION_DOUBLES,
		GenderCategory:        enums.CATEGORY_GENDER_MIXED,
		Level:                 enums.CATEGORY_LEVEL_BEGINNER,
		LinkWAGroup:           "https://chat.whatsapp.com/JUup3ZHCEMd1eO2WifCsI5",
		Season:                1,
		RegistrationFee:       400_000,
		LinkWARatingAdmin:     "628159239794",
	},
}

func ResetLeagues(db *gorm.DB) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if !utils.IsDevEnv() {
			return errors.New("reset leagues can only be done in dev environment")
		}

		tx := db.WithContext(c.Context).Begin()

		err := tx.
			Unscoped().
			Delete(&entities.LeaguePlayer{}, "deleted_at IS NULL").
			Error
		if err != nil {
			return fmt.Errorf("tx.Select(clause.Associations).Unscoped().Delete(&entities.LeaguePlayer{}).Error: %w", err)
		}

		err = tx.
			Unscoped().
			Delete(&entities.LeagueTeam{}, "deleted_at IS NULL").
			Error
		if err != nil {
			return fmt.Errorf("tx.Select(clause.Associations).Unscoped().Delete(&entities.LeagueTeam{}).Error: %w", err)
		}

		err = tx.
			Unscoped().
			Delete(&entities.LeagueRegistration{}, "deleted_at IS NULL").
			Error
		if err != nil {
			return fmt.Errorf("tx.Select(clause.Associations).Unscoped().Delete(&entities.LeagueTeam{}).Error: %w", err)
		}

		err = tx.
			Unscoped().
			Delete(&entities.League{}, "deleted_at IS NULL").
			Error
		if err != nil {
			return fmt.Errorf("tx.Select(clause.Associations).Unscoped().Delete(&entities.League{}).Error: %w", err)
		}

		err = tx.Create(DEFAULT_LEAGUES).Error
		if err != nil {
			return fmt.Errorf("tx.Create(DEFAULT_LEAGUES).Error: %w", err)
		}

		err = tx.Commit().Error
		if err != nil {
			return err
		}

		return nil
	}
}
