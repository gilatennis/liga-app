package repositories

import (
	"context"
	"liga-app/entities"

	"github.com/google/uuid"
)

type IMatchSetRepository interface {
	SaveMatchSet(ctx context.Context, userUUID uuid.UUID, matchSet *entities.MatchSet) (*uuid.UUID, error)
}
