package entities

import "github.com/google/uuid"

type Admin struct {
	Base
	UserUUID uuid.UUID `gorm:"type:CHAR(36)"`
	User     User      `gorm:"ForeignKey:UserUUID;References:UUID"`
}
