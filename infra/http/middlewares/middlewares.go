package middlewares

import (
	"liga-app/config"
	"liga-app/interfaces"

	// jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Middlewares struct {
	// JwtMiddleware     *jwt.GinJWTMiddleware
	SessionMiddleware gin.HandlerFunc
	HistoryMiddleware gin.HandlerFunc
}

func SetupMiddlewares(config *config.Config, db *gorm.DB, fileStorage interfaces.IFileStorage) (*Middlewares, error) {
	/*
		jwtMiddleware, err := NewJwtMiddleware(config, db)
		if err != nil {
			return nil, fmt.Errorf("error creating jwt middleware: %w", err)
		}
	*/

	sessionMiddleware := NewSessionMiddleware(config, db, fileStorage)
	historyMiddleWare := NewHistoryMiddleware()

	return &Middlewares{
		// JwtMiddleware:     jwtMiddleware,
		SessionMiddleware: sessionMiddleware,
		HistoryMiddleware: historyMiddleWare,
	}, nil
}
