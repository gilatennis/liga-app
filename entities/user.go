package entities

import "liga-app/enums"

type User struct {
	Base
	Email       string
	Name        string
	Phone       string
	Province    string
	Gender      enums.GENDER
	LoginMethod enums.LOGIN_METHOD
	Ratings     []UserRating `gorm:"ForeignKey:UserUUID"`
}
