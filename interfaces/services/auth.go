package services

import (
	"context"
	"liga-app/entities"
	"liga-app/enums"

	"github.com/google/uuid"
)

type IRegisterPayload interface {
	GetName() string
	GetEmail() string
	GetPassword() string
}

type IAuthService interface {
	GetMe(ctx context.Context, userUUID uuid.UUID) (*entities.User, error)

	RegisterWithMagicLink(ctx context.Context, email string, redirect *string) error
	CompleteMagicLinkRegistration(ctx context.Context, email string, token string) (*entities.User, error) // returns redirect url, error

	LoginWithMagicLink(ctx context.Context, email string, redirect *string) error
	CompleteMagicLinkLogin(ctx context.Context, email string, token string) (*entities.User, error) // returns redirect url, error

	LoginWithOAuth(ctx context.Context, email string, loginMethod enums.LOGIN_METHOD) (*entities.User, error)
}
