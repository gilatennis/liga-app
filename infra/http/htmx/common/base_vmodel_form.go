package common

import (
	"fmt"
	"net/mail"
	"reflect"
	"regexp"
	"time"
)

type BaseVModelForm struct {
	FormErrors  map[string]string
	ServerError string
	Test        string `form:"test" validation:"test" label:"test"`
}

func (b *BaseVModelForm) Validate(parentModel any) bool {
	b.FormErrors = make(map[string]string)
	isValid := true

	modelTypes := reflect.TypeOf(parentModel)
	modelValues := reflect.ValueOf(parentModel)

	for i := 0; i < modelTypes.NumField(); i++ {
		field := modelTypes.Field(i)
		fieldValue := modelValues.Field(i)
		name := field.Tag.Get("form")
		validation := field.Tag.Get("validation")
		label := field.Tag.Get("label")

		if validation == "" {
			continue
		}

		if validation == "required" {
			if fieldValue.String() == "" {
				b.FormErrors[name] = fmt.Sprintf("%s is required", label)
				isValid = false
			}

			fieldType := fieldValue.Type().String()
			if fieldType == "time.Time" {
				if fieldValue.Interface().(time.Time).IsZero() {
					b.FormErrors[name] = fmt.Sprintf("%s is required", label)
					isValid = false
				}
			}
		}

		if validation == "email" {
			// TODO: Check if valid email
			_, err := mail.ParseAddress(fieldValue.String())
			if err != nil {
				b.FormErrors[name] = fmt.Sprintf("%s must be a valid email", label)
				isValid = false
			}
		}

		if validation == "phone" {
			r, err := regexp.Compile(`^\+?\d{7,14}$`)
			if err != nil {
				b.ServerError = err.Error()
				isValid = false
			}

			isValid = r.MatchString(fieldValue.String())
			if !isValid {
				b.FormErrors[name] = fmt.Sprintf("%s invalid", label)
			}
		}

		// TODO: add other basic validations here
	}

	return isValid
}
