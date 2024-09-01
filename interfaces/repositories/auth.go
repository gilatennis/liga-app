package repositories

import (
	"context"
	"liga-app/entities"

	"github.com/google/uuid"
)

type IAuthRepository interface {
	CreateMagicLinkRegistration(ctx context.Context, reg *entities.MagicLinkRegistration) (*uuid.UUID, error)
	UpdateMagicLinkRegistration(ctx context.Context, reg *entities.MagicLinkRegistration) (*uuid.UUID, error)
	GetMagicLinkRegistrationByEmail(ctx context.Context, email string) (*entities.MagicLinkRegistration, error)

	CreateMagicLinkLogin(ctx context.Context, login *entities.MagicLinkLogin) (*uuid.UUID, error)
	UpdateMagicLinkLogin(ctx context.Context, login *entities.MagicLinkLogin) (*uuid.UUID, error)
	GetMagicLinkLoginByToken(ctx context.Context, token string) (*entities.MagicLinkLogin, error)
}
