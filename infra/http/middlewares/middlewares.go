package middlewares

import (
	"liga-app/config"

	// jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Middlewares struct {
	// JwtMiddleware     *jwt.GinJWTMiddleware
	SessionMiddleware gin.HandlerFunc
}

func SetupMiddlewares(config *config.Config, db *gorm.DB) (*Middlewares, error) {
	/*
		jwtMiddleware, err := NewJwtMiddleware(config, db)
		if err != nil {
			return nil, fmt.Errorf("error creating jwt middleware: %w", err)
		}
	*/

	sessionMiddleware := NewSessionMiddleware(config, db)

	return &Middlewares{
		// JwtMiddleware:     jwtMiddleware,
		SessionMiddleware: sessionMiddleware,
	}, nil
}
