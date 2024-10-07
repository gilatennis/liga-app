package entities

import (
	"fmt"
	"liga-app/enums"
	"sort"
)

type User struct {
	Base
	Email             string
	Name              string
	Phone             string
	Province          string
	Gender            enums.GENDER
	LoginMethod       enums.LOGIN_METHOD
	Ratings           []UserRating `gorm:"ForeignKey:UserUUID"`
	ProfilePicture    string
	ProfilePictureURL string `gorm:"-:all"`
	IsTestUser        bool
}

func (u *User) GetLatestRatingString() string {
	if len(u.Ratings) == 0 {
		return "-"
	}

	sort.SliceStable(u.Ratings, func(i, j int) bool {
		return u.Ratings[i].NTRPRating > u.Ratings[j].NTRPRating
	})

	return fmt.Sprintf("NTRP %.1f", u.Ratings[0].NTRPRating)
}
