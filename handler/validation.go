package handler

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var userValidation map[string]string = map[string]string{
	"name":     "required,min=3,max=60",
	"phone":    "required,min=10,max=13",
	"password": "required,min=6,max=64",
}

var errorMapp map[string]string = map[string]string{
	"required": "value is required",
	"min":      "should be greater than",
	"max":      "should be less than",
}

func UserValidation(user interface{}) (bool, []string) {
	rt := reflect.TypeOf(user)
	values := reflect.ValueOf(user)
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	var errorMessages []string
	isValid := true
	message := ""

	for i := 0; i < values.NumField(); i++ {
		isValid = true
		var (
			min int
			max int
		)

		jsonTag := strings.Split(rt.Field(i).Tag.Get("json"), ",")[0]

		fmt.Sscanf(userValidation[jsonTag], "required,min=%d,max=%d", &min, &max)

		val := values.Field(i).Interface()
		if ok := strings.Contains(userValidation[jsonTag], "required"); ok {
			if res, ok := val.(string); ok && res == "" {
				isValid = false
				message = fmt.Sprintf("%s %s", jsonTag, errorMapp["required"])
				errorMessages = append(errorMessages, message)
			}
		}

		if isValid {
			if ok := strings.Contains(userValidation[jsonTag], "max"); ok {
				if len(val.(string)) > max {
					isValid = false
					message = fmt.Sprintf("%s %s %d %s", jsonTag, errorMapp["max"], max, "character")
					errorMessages = append(errorMessages, message)

				}
			}

			if ok := strings.Contains(userValidation[jsonTag], "min"); ok {
				if len(val.(string)) < min {
					isValid = false
					message += fmt.Sprintf("%s %s %d %s", jsonTag, errorMapp["min"], min, "character")
					errorMessages = append(errorMessages, message)
				}
			}
		}

		switch jsonTag {
		case "phone":
			isValid, errorMessages = PhoneValidation(val.(string), errorMessages)
		case "password":
			isValid, errorMessages = PasswordValidation(val.(string), errorMessages)
		}

	}

	return isValid, errorMessages
}

func PhoneValidation(phone string, errList []string) (bool, []string) {
	var errorMessages []string
	isValid := true
	message := ""
	pattern := regexp.MustCompile(`^\+62\d{8,11}$`)
	matched := pattern.MatchString(phone)

	if !matched {
		isValid = false
		message = fmt.Sprintf(`%s`, "Phone must start with +62")
		errorMessages = append(errorMessages, message)
	}

	errList = append(errList, errorMessages...)
	return isValid, errList

}

func PasswordValidation(password string, errList []string) (bool, []string) {
	var errorMessages []string
	isValid := true
	message := ""
	pattern := regexp.MustCompile(`[A-Z]`)
	containsCapital := pattern.MatchString(password)
	if !containsCapital {
		isValid = false
		message = fmt.Sprintf(`%s`, "Password contains at least 1 Capital")
		errorMessages = append(errorMessages, message)
	}
	pattern = regexp.MustCompile(`\d`)
	containsNumber := pattern.MatchString(password)

	if !containsNumber {
		isValid = false
		message = "Password contains at least 1 number"
		errorMessages = append(errorMessages, message)
	}
	pattern = regexp.MustCompile(`\W`)
	containsSpecialChar := pattern.MatchString(password)
	if !containsSpecialChar {
		isValid = false
		message = "Password contains at least 1 spesial character"
		errorMessages = append(errorMessages, message)
	}

	errList = append(errList, errorMessages...)
	return isValid, errList
}
