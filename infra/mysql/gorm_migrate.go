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
	)

	if errMigrate != nil {
		return fmt.Errorf("[mysql.Migrate]: error migrating: %w", errMigrate)
	}

	return nil
}
