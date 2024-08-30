package entities

import "time"

type MagicLinkRegistration struct {
	Base
	Email          string
	Token          string
	TokenExpiredAt time.Time
	VerifiedAt     *time.Time
}

type MagicLinkLogin struct {
	Base
	Email          string
	Token          string
	TokenExpiredAt time.Time
	VerifiedAt     *time.Time
}
