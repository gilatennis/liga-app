package middlewares

import (
	"context"
	"errors"
	"fmt"
	"liga-app/config"
	"liga-app/entities"
	"liga-app/interfaces"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func NewSessionMiddleware(
	cfg *config.Config,
	db *gorm.DB,
	fileStorage interfaces.IFileStorage,
) gin.HandlerFunc {
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

		fullPath := c.FullPath()
		if strings.HasPrefix("/admin/dashboard", fullPath) {
			isAdmin, err := sessionIsAdmin(c, db, userUUID)
			if err != nil {
				fmt.Printf("sessionIsAdmin(c, db, userUUID=%s) : %+v \n", userUUIDStr, err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}

			if !isAdmin {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
		}

		user, err := sessionGetUserByUUID(c, cfg, db, fileStorage, userUUID)
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

func sessionGetUserByUUID(
	ctx context.Context,
	cfg *config.Config,
	db *gorm.DB,
	fileStorage interfaces.IFileStorage,
	userUUID uuid.UUID,
) (*entities.User, error) {
	var user entities.User
	err := db.Where("uuid = ?", userUUID).
		Preload("Ratings").
		First(&user).
		Error
	if err != nil {
		return nil, err
	}

	if user.ProfilePicture != "" {
		user.ProfilePictureURL, err = fileStorage.GetFileURL(ctx, cfg.Minio.Bucket, user.ProfilePicture)
		if err != nil {
			fmt.Printf("fileStorage.GetFileURL for user.ProfilePicture=%s failed: %v \n", user.ProfilePicture, err)
		}
	}

	return &user, nil
}

func sessionIsAdmin(
	ctx context.Context,
	db *gorm.DB,
	userUUID uuid.UUID,
) (bool, error) {
	adm := entities.Admin{}

	err := db.WithContext(ctx).Where("user_uuid = ?", userUUID).First(&adm).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("s.db.Where(user_uuid = ?, userUUID=%s): %w", userUUID, err)
	}

	return true, nil
}
