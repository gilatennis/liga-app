package users

import (
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return nil
	}
}
