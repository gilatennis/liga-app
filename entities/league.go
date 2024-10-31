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
	RegistrationFee       float64
	LinkWARatingAdmin     string
	PlayerCount           int64     `gorm:"-:all"`
	TeamCount             int64     `gorm:"-:all"`
	Fixtures              []Fixture `gorm:"ForeignKey:LeagueUUID"`
}

func (l *League) IsNTRPRatingTooHigh(rating float32) bool {
	if l.Level == enums.CATEGORY_LEVEL_BEGINNER {
		return rating > 3.5
	}

	if l.Level == enums.CATEGORY_LEVEL_INTERMEDIATE {
		return rating > 4
	}

	return false
}

func (l *League) IsNTRPRatingTooLow(rating float32) bool {
	if l.Level == enums.CATEGORY_LEVEL_ADVANCED {
		return rating < 4.5
	}

	if l.Level == enums.CATEGORY_LEVEL_INTERMEDIATE {
		return rating < 3.5
	}

	return false
}
