package entities

import (
	"liga-app/enums"
	"time"
)

type League struct {
	Base
	Name                  string
	WilayahIDProvince     string                       `gorm:"type:varchar(24)"`
	WilayahIDCity         string                       `gorm:"type:varchar(24)"`
	WilayahIDSubdistrict  string                       `gorm:"type:varchar(24)"`
	WilayahIDVillage      string                       `gorm:"type:varchar(24)"`
	Wilayah               *Wilayah                     `gorm:"ForeignKey:WilayahIDProvince,WilayahIDCity,WilayahIDSubdistrict,WilayahIDVillage;References:IDProvince,IDCity,IDSubdistrict,IDVillage"`
	ParticipationCategory enums.CATEGORY_PARTICIPATION `gorm:"type:enum('SINGLES', 'DOUBLES', 'TEAMS')"`
	GenderCategory        enums.CATEGORY_GENDER        `gorm:"type:enum('MALE', 'FEMALE', 'MIXED')"`
	Level                 enums.CATEGORY_LEVEL         `gorm:"type:enum('BEGINNER', 'INTERMEDIATE', 'ADVANCED')"`
	LinkRuleBook          string
	LinkWAGroup           string
	RegistrationStartAt   *time.Time
	RegistrationEndAt     *time.Time
	LeagueStartAt         *time.Time
	LeagueEndAt           *time.Time
	Season                int
}
