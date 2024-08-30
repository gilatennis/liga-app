package repositories

import (
	"context"
	"liga-app/entities"

	"github.com/google/uuid"
)

type IAuthRepository interface {
	// CreateEmailRegistration(ctx context.Context, reg *entities.EmailRegistration) error
	// UpdateEmailRegistration(ctx context.Context, reg *entities.EmailRegistration) error
	// GetEmailRegistrationByEmail(ctx context.Context, email string) (*entities.EmailRegistration, error)
	// DeleteEmailRegistrationByEmail(ctx context.Context, email string) error
	// CreateResetPasswordRequest(ctx context.Context, req *entities.ResetPasswordRequest) (uuid.UUID, error)
	// GetResetPasswordRequestByToken(ctx context.Context, token string) (*entities.ResetPasswordRequest, error)

	CreateMagicLinkRegistration(ctx context.Context, reg *entities.MagicLinkRegistration) (*uuid.UUID, error)
	UpdateMagicLinkRegistration(ctx context.Context, reg *entities.MagicLinkRegistration) (*uuid.UUID, error)
	GetMagicLinkRegistrationByEmail(ctx context.Context, email string) (*entities.MagicLinkRegistration, error)

	CreateMagicLinkLogin(ctx context.Context, login *entities.MagicLinkLogin) (*uuid.UUID, error)
	UpdateMagicLinkLogin(ctx context.Context, login *entities.MagicLinkLogin) (*uuid.UUID, error)
	GetMagicLinkLoginByToken(ctx context.Context, token string) (*entities.MagicLinkLogin, error)
}
