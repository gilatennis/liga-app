package entities

import "liga-app/enums"

type User struct {
	Base
	Email       string
	Name        string
	Phone       string
	Gender      enums.GENDER
	LoginMethod enums.LOGIN_METHOD
}
