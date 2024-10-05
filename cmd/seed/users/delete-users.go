package users

import (
	"errors"
	"fmt"
	"liga-app/entities"
	"liga-app/utils"

	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func DeleteUsers(db *gorm.DB) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if !utils.IsDevEnv() {
			return errors.New("reset leagues can only be done in dev environment")
		}

		err := db.WithContext(c.Context).
			Delete(&entities.User{}, "email != 'kevin.janada@gmail.com'").
			Error
		if err != nil {
			return fmt.Errorf("db.WithContext(c.Context).Delete(&entities.User{}, \"email != 'kevin.janada@gmail.com'\").Error: %w", err)
		}

		return nil
	}
}
