package user_ratings

import (
	"fmt"
	"liga-app/entities"
	"liga-app/infra/gsheet"

	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func UpdateUserRatings(db *gorm.DB, gsheetAPI *gsheet.API) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		userRatings, err := gsheetAPI.UserRatings.GetUserRatings()
		if err != nil {
			return fmt.Errorf("gsheetAPI.UserRatings.GetUserRatings() : %w", err)
		}

		for _, ur := range userRatings {
			fmt.Printf("%s %.1f \n", ur.UserUUID, ur.NTRPRating)

			err := db.Model(&entities.UserRating{}).
				Where(entities.UserRating{UserUUID: ur.UserUUID}).
				FirstOrCreate(&ur).
				Error
			if err != nil {
				return fmt.Errorf("db.Model(&entities.UserRating{})...FirstOrCreate(&ur=%+v): %w", ur, err)
			}
		}

		return nil
	}
}
