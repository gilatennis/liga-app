package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	UUID      uuid.UUID `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uuid.UUID  `gorm:"type:char(36)"`
	UpdatedBy uuid.UUID  `gorm:"type:char(36)"`
	DeletedAt *time.Time `gorm:"index"`
	DeletedBy *uuid.UUID `gorm:"type:char(36)"`
}

func (w *Base) BeforeCreate(tx *gorm.DB) (err error) {
	w.UUID = uuid.New()
	w.CreatedAt = time.Now()
	w.UpdatedAt = time.Now()
	return
}

func (w *Base) SetCreatedBy(user *User) {
	w.CreatedBy = user.UUID
}

func (w *Base) SetUpdatedBy(user *User) {
	w.UpdatedBy = user.UUID
}

func (w *Base) SetDeletedBy(user *User) {
	w.DeletedBy = &user.UUID
}

func (w *Base) GetUUID() *uuid.UUID {
	return &w.UUID
}
