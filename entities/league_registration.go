package entities

import (
	"liga-app/enums"

	"github.com/google/uuid"
)

type LeagueRegistration struct {
	Base
	LeagueUUID    uuid.UUID `gorm:"type:char(36)"`
	League        League    `gorm:"ForeignKey:LeagueUUID;References:UUID"`
	UserUUID      uuid.UUID `gorm:"type:char(36)"`
	User          User      `gorm:"ForeignKey:UserUUID;References:UUID"`
	PaymentAmount float64
	PaymentStatus enums.PAYMENT_STATUS `gorm:"type:enum('UNPAID', 'PAID')"`
	PaymentID     string
}
