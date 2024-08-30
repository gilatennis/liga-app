package mysql

import (
	"fmt"
	"liga-app/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DSN=y675ro26w9547806y0c5:pscale_pw_XY8ogW8Nm4msFyDvdS8Is0g8qhpaUWnLEsGEZ6KzWu5@tcp(aws.connect.psdb.cloud)/arise?tls=true&parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci

func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
	tls := "true"
	if !cfg.DB.TLS {
		tls = "false"
	}

	dsn := cfg.DB.User + ":" + cfg.DB.Password +
		"@tcp(" + cfg.DB.Host + ":" + cfg.DB.Port + ")/" + cfg.DB.Name +
		"?tls=" + tls +
		"&parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("[infra.mysql.db.NewDB] gorm.Open(): failed to open database: %w", err)
	}

	return db, nil
}
