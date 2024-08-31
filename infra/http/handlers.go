package http

import (
	"fmt"
	"liga-app/config"
	"liga-app/infra/http/htmx/auth"
	"liga-app/infra/http/htmx/dashboard"
	"liga-app/infra/http/htmx/index"
	"liga-app/infra/http/json/duitku"
	"liga-app/infra/minio"
	"liga-app/third_party/mailtrap"

	domainAuth "liga-app/domain/auth"
	domainUserInfo "liga-app/domain/user_info"
	dbAuth "liga-app/infra/mysql/auth"
	dbUser "liga-app/infra/mysql/user"
	dbWilayah "liga-app/infra/mysql/wilayah"

	// "liga-app/mixins"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Handlers struct {
	HTMXHandlers *HTMXHandlers
	JSONHandlers *JSONHandlers
}

type HTMXHandlers struct {
	IndexHandler     *index.Handler
	AuthHandler      *auth.Handler
	DashboardHandler *dashboard.Handler
}

type JSONHandlers struct {
	DuitkuHandler *duitku.Handler
}

func setupHandlers(config *config.Config, db *gorm.DB, logger *zerolog.Logger) (*Handlers, error) {
	mail := mailtrap.NewMailTrap(config)
	minioClient, err := minio.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("http.server.setupHandlers(): failed to create new minio client: %w", err)
	}
	_ = mail
	_ = minioClient

	// fileUploader := mixins.NewFileUploader(minioClient)

	/*
		accessControlProvider, err := casbin.NewProvider(db)
		if err != nil {
			return nil, fmt.Errorf("casbin.NewProvider(db) %w", err)
		}
		accessControlStore := accessControlDB.NewStore(db)
	*/

	authStore := dbAuth.NewStore(db)
	userStore := dbUser.NewStore(db)
	wilayahStore := dbWilayah.NewStore(db)

	authService := domainAuth.NewService(
		config,
		authStore,
		userStore,
		mail,
	)
	userInfoService := domainUserInfo.NewService(wilayahStore, userStore)

	HTMXHandlers := HTMXHandlers{
		AuthHandler:      auth.NewHandler(config, authService),
		IndexHandler:     index.NewHandler(),
		DashboardHandler: dashboard.NewHandler(userInfoService),
	}

	JSONHandlers := JSONHandlers{
		DuitkuHandler: duitku.NewHandler(),
	}

	return &Handlers{
		HTMXHandlers: &HTMXHandlers,
		JSONHandlers: &JSONHandlers,
	}, nil
}
