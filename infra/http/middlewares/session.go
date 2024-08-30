package middlewares

import (
	"fmt"
	"liga-app/config"
	"liga-app/entities"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func NewSessionMiddleware(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userUUIDAny := session.Get("user_uuid")
		userUUIDStr, ok := session.Get("user_uuid").(string)

		if userUUIDAny == nil || !ok {
			// Abort the request with the appropriate error code
			fmt.Println("user_uuid not found in session")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userUUID, err := uuid.Parse(userUUIDStr)
		if err != nil {
			fmt.Println("error parsing user_uuid from session. user_uuid: ", userUUIDStr)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		user, err := sessionGetUserByUUID(db, userUUID)
		if err != nil {
			err = fmt.Errorf("error getting user by user_uuid from session. user_uuid: %s, err: %w", userUUIDStr, err)
			fmt.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("user_uuid", userUUIDStr)
		c.Set("user", user)
		// Continue down the chain to handler etc
		c.Next()
	}
}

func sessionGetUserByUUID(db *gorm.DB, userUUID uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := db.Where("uuid = ?", userUUID).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
