package auth

import (
	"context"
	"errors"
	"fmt"
	"liga-app/config"
	domainErrors "liga-app/domain/errors"
	"liga-app/entities"
	"liga-app/enums"
	"liga-app/interfaces"
	"liga-app/interfaces/repositories"
	"time"

	"github.com/dchest/uniuri"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var logTag = "[domain.auth.Service]"

type Service struct {
	config      *config.Config
	authRepo    repositories.IAuthRepository
	userRepo    repositories.IUserRepository
	emailSender interfaces.IEmailSender
}

func NewService(
	config *config.Config,
	authRepo repositories.IAuthRepository,
	userRepo repositories.IUserRepository,
	emailSender interfaces.IEmailSender,
) *Service {
	return &Service{
		config:      config,
		authRepo:    authRepo,
		userRepo:    userRepo,
		emailSender: emailSender,
	}
}

func (s *Service) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *Service) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *Service) GenerateToken() string {
	return uniuri.NewLen(20)
}

func (s *Service) RegisterWithMagicLink(ctx context.Context, email string, redirect *string) error {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	var errNotFound *domainErrors.ErrNotFound
	if err != nil && !errors.As(err, &errNotFound) {
		return fmt.Errorf("%s: %w", logTag, err)
	}

	if user != nil {
		err = &domainErrors.ErrConflict{Message: fmt.Sprintf("user with email %s already exists", email)}
		return fmt.Errorf("%s: %w", logTag, err)
	}

	magicLinkReg, err := s.authRepo.GetMagicLinkRegistrationByEmail(ctx, email)
	if err != nil && !errors.As(err, &errNotFound) {
		return fmt.Errorf("%s: %w", logTag, err)
	}

	if magicLinkReg != nil {
		magicLinkReg.Token = s.GenerateToken()
		magicLinkReg.TokenExpiredAt = time.Now().Add(time.Hour * 24)
		_, err = s.authRepo.UpdateMagicLinkRegistration(ctx, magicLinkReg)
		if err != nil {
			return fmt.Errorf("%s: %w", logTag, err)
		}
	} else {
		magicLinkReg = &entities.MagicLinkRegistration{
			Email:          email,
			Token:          s.GenerateToken(),
			TokenExpiredAt: time.Now().Add(time.Hour * 24),
		}

		_, err = s.authRepo.CreateMagicLinkRegistration(ctx, magicLinkReg)
		if err != nil {
			return fmt.Errorf("%s: %w", logTag, err)
		}
	}

	magicLinkEmailContent := s.CreateRegisterMagicLinkEmailContent(email, magicLinkReg.Token, redirect)

	from := s.config.Email.EmailSender
	to := email
	subject := "Magic Link Registration"
	content := magicLinkEmailContent

	err = s.emailSender.SendEmail(ctx, from, to, subject, content)
	if err != nil {
		return fmt.Errorf("%s: %w", logTag, err)
	}

	return nil
}

func (s *Service) CompleteMagicLinkRegistration(ctx context.Context, email string, token string) (*entities.User, error) {
	magLinkReg, err := s.authRepo.GetMagicLinkRegistrationByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", logTag, err)
	}

	if magLinkReg.VerifiedAt != nil {
		err = &domainErrors.ErrBadRequest{Message: "magic link already used"}
		return nil, fmt.Errorf("%s: %w", logTag, err)
	}

	if magLinkReg.Token != token {
		err = &domainErrors.ErrUnauthorized{Message: "invalid magic link"}
		return nil, fmt.Errorf("%s: %w", logTag, err)
	}

	if time.Now().After(magLinkReg.TokenExpiredAt) {
		err = &domainErrors.ErrUnauthorized{Message: "magic link expired"}
		return nil, fmt.Errorf("%s: %w", logTag, err)
	}

	user := &entities.User{
		Email:       email,
		LoginMethod: enums.LOGIN_METHOD_EMAIL,
	}

	userUUID, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", logTag, err)
	}

	user.UUID = *userUUID

	now := time.Now()
	magLinkReg.VerifiedAt = &now

	_, err = s.authRepo.UpdateMagicLinkRegistration(ctx, magLinkReg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", logTag, err)
	}

	return user, nil
}

func (s *Service) LoginWithMagicLink(ctx context.Context, email string, redirect *string) error {
	_, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("%s: %w", logTag, err)
	}

	magicLinkToken := s.GenerateToken()
	tokenExpiredAt := time.Now().Add(time.Hour * 24)

	magicLinkLogin := &entities.MagicLinkLogin{
		Email:          email,
		Token:          magicLinkToken,
		TokenExpiredAt: tokenExpiredAt,
	}

	_, err = s.authRepo.CreateMagicLinkLogin(ctx, magicLinkLogin)
	if err != nil {
		return fmt.Errorf("%s: %w", logTag, err)
	}

	magicLinkEmailContent := s.CreateLoginMagicLinkEmailContent(email, magicLinkToken, redirect)

	from := s.config.Email.EmailSender
	to := email
	subject := "Magic Link Login"
	content := magicLinkEmailContent

	err = s.emailSender.SendEmail(ctx, from, to, subject, content)
	if err != nil {
		return fmt.Errorf("%s: %w", logTag, err)
	}

	return nil
}

func (s *Service) CompleteMagicLinkLogin(ctx context.Context, email string, token string) (*entities.User, error) {
	magicLinkLogin, err := s.authRepo.GetMagicLinkLoginByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", logTag, err)
	}

	if magicLinkLogin.VerifiedAt != nil {
		err = &domainErrors.ErrBadRequest{Message: "magic link already used"}
		return nil, fmt.Errorf("%s: %w", logTag, err)
	}

	if time.Now().After(magicLinkLogin.TokenExpiredAt) {
		err = &domainErrors.ErrUnauthorized{Message: "magic link expired"}
		return nil, fmt.Errorf("%s: %w", logTag, err)
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", logTag, err)
	}

	now := time.Now()
	magicLinkLogin.VerifiedAt = &now

	_, err = s.authRepo.UpdateMagicLinkLogin(ctx, magicLinkLogin)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", logTag, err)
	}

	return user, nil
}

func (s *Service) LoginWithOAuth(ctx context.Context, email string, loginMethod enums.LOGIN_METHOD) (*entities.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	var errNotFound *domainErrors.ErrNotFound
	if err != nil && !errors.As(err, &errNotFound) {
		return nil, fmt.Errorf("%s: %w", logTag, err)
	}

	if user == nil {
		user = &entities.User{
			Email:       email,
			LoginMethod: loginMethod,
		}
		userUUID, err := s.userRepo.CreateUser(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", logTag, err)
		}
		user.UUID = *userUUID
		return user, nil
	}

	user.LoginMethod = loginMethod
	_, err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", logTag, err)
	}

	return user, nil
}

func (s *Service) CreateVerificationEmailContent(email string, verificationToken string) string {
	serverBaseUrl := s.config.App.BaseURL

	text := "Please click this link to verify your email:" +
		serverBaseUrl + "/auth/verify-email?email=" + email +
		"&token=" + verificationToken

	return text
}

func (s *Service) CreateRegisterMagicLinkEmailContent(email string, magicLinkToken string, redirect *string) string {
	serverBaseUrl := s.config.App.BaseURL

	text := "Please click this link to complete registration:" +
		serverBaseUrl + "/auth/register/magic-link?email=" + email +
		"&token=" + magicLinkToken

	if redirect != nil {
		text += "&redirect=" + *redirect
	}

	return text
}

func (s *Service) CreateLoginMagicLinkEmailContent(email string, magicLinkToken string, redirect *string) string {
	serverBaseUrl := s.config.App.BaseURL

	text := "Please click this link to login:" +
		serverBaseUrl + "/auth/login/magic-link?email=" + email +
		"&token=" + magicLinkToken

	if redirect != nil {
		text += "&redirect=" + *redirect
	}

	return text
}

func (s *Service) CreateResetPasswordEmailContent(token string) string {
	baseURL := s.config.App.WebAppBaseURL

	text := "Please click this link to reset your password: " +
		baseURL + "/forgot-password/update?token=" + token

	return text
}

/*
func (s *Service) Register(ctx context.Context, payload services.IRegisterPayload) error {
	name := payload.GetName()
	email := payload.GetEmail()
	password := payload.GetPassword()
	hashedPassword := ""

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	var errNotFound *domainErrors.ErrNotFound
	if err != nil && !errors.As(errors.Unwrap(err), &errNotFound) {
		return err
	}

	emailReg, err := s.authRepo.GetEmailRegistrationByEmail(ctx, email)
	if err != nil && !errors.As(err, &errNotFound) {
		return err
	}

	if user != nil && emailReg != nil && emailReg.IsVerified {
		return &domainErrors.ErrConflict{Message: fmt.Sprintf("user with email %s already exists", email)}
	}

	if user != nil && emailReg != nil && !emailReg.IsVerified {
		return &domainErrors.ErrConflict{Message: fmt.Sprintf("the email %s is registered, but email is not verified", email)}
	}

	hashedPassword, err = s.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user = &entities.User{
		Name:        name,
		Email:       email,
		Password:    hashedPassword,
		LoginMethod: enums.LoginMethodEmail,
	}
	_, err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	verificationToken := s.GenerateToken()
	expiresAt := time.Now().Add(time.Hour * 24)
	if emailReg == nil {
		emailReg = &entities.EmailRegistration{
			Name:           name,
			Email:          email,
			Password:       hashedPassword,
			Token:          verificationToken,
			TokenExpiredAt: expiresAt,
		}
		err := s.authRepo.CreateEmailRegistration(ctx, emailReg)
		if err != nil {
			return err
		}
	} else {
		emailReg.Password = hashedPassword
		emailReg.Token = verificationToken
		emailReg.TokenExpiredAt = expiresAt
		err = s.authRepo.UpdateEmailRegistration(ctx, emailReg)
		if err != nil {
			return err
		}
	}

	verifEmailContent := s.CreateVerificationEmailContent(email, verificationToken)

	from := s.config.Email.EmailSender
	to := email
	subject := "Email Verification"
	content := verifEmailContent

	err = s.emailSender.SendEmail(ctx, from, to, subject, content)
	if err != nil {
		return err
	}

	return nil
}
*/

/*
func (s *Service) VerifyEmail(ctx context.Context, email string, verificationToken string) (string, error) {
	emailReg, err := s.authRepo.GetEmailRegistrationByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if emailReg.IsVerified {
		return "", &domainErrors.ErrBadRequest{Message: "email already verified"}
	}

	if emailReg.Token != verificationToken {
		return "", &domainErrors.ErrUnauthorized{Message: "invalid verification token"}
	}

	if time.Now().After(emailReg.TokenExpiredAt) {
		return "", &domainErrors.ErrUnauthorized{Message: "verification token expired"}
	}

	emailReg.IsVerified = true
	err = s.authRepo.UpdateEmailRegistration(ctx, emailReg)
	if err != nil {
		return "", err
	}

	// TODO: if request sender is web app, redirect to web app login page
	// if request sender is mobile app, return deeplink to mobile app login page

	redirectUrl := s.config.App.WebAppBaseURL + "/login" + "?email=" + email

	return redirectUrl, nil
}
*/

func (s *Service) GetMe(ctx context.Context, userUUID uuid.UUID) (*entities.User, error) {
	user, err := s.userRepo.GetUserByUUID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

/*
func (s *Service) ResendVerificationEmail(ctx context.Context, email string) error {
	emailReg, err := s.authRepo.GetEmailRegistrationByEmail(ctx, email)
	if err != nil {
		return err
	}

	if emailReg.IsVerified {
		return &domainErrors.ErrBadRequest{Message: "email already verified"}
	}

	verificationToken := s.GenerateToken()
	expiresAt := time.Now().Add(time.Hour * 24)
	emailReg.Token = verificationToken
	emailReg.TokenExpiredAt = expiresAt
	err = s.authRepo.UpdateEmailRegistration(ctx, emailReg)
	if err != nil {
		return err
	}

	verifEmailContent := s.CreateVerificationEmailContent(email, emailReg.Token)

	from := s.config.Email.EmailSender
	to := email
	subject := "Email Verification"
	content := verifEmailContent

	err = s.emailSender.SendEmail(ctx, from, to, subject, content)
	if err != nil {
		return err
	}

	return nil
}
*/

/*
func (s *Service) CreateResetPasswordRequest(ctx context.Context, email string) error {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	resetPasswordToken := uniuri.NewLen(20)
	tokenExpiresAt := time.Now().Add(time.Hour * 24)

	resetPasswordReq := &entities.ResetPasswordRequest{
		UserUUID:       user.UUID,
		User:           *user,
		Token:          resetPasswordToken,
		TokenExpiredAt: tokenExpiresAt,
	}

	_, err = s.authRepo.CreateResetPasswordRequest(ctx, resetPasswordReq)
	if err != nil {
		return err
	}

	emailContent := s.CreateResetPasswordEmailContent(resetPasswordToken)
	from := s.config.Email.EmailSender
	to := user.Email
	subject := "Reset Password"
	content := emailContent

	err = s.emailSender.SendEmail(ctx, from, to, subject, content)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
*/

/*
func (s *Service) UpdatePassword(ctx context.Context, token string, password string) error {
	req, err := s.authRepo.GetResetPasswordRequestByToken(ctx, token)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetUserByUUID(ctx, req.UserUUID)
	if err != nil {
		return err
	}

	// TODO: password validation

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	_, err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
*/
