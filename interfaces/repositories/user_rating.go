package repositories

import (
	"context"
	"liga-app/entities"

	"github.com/google/uuid"
)

type IUserRatingRepository interface {
	Save(ctx context.Context, userRating *entities.UserRating) (*uuid.UUID, error)
}
