package handler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/go-playground/validator/v10"
)

var _validator *validator.Validate

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

func UserValidation(user generated.RegisterUserRequest) (bool, map[string][]string) {
	var errList map[string][]string
	isValid := true
	message := ""
	fmt.Println(user)
	for i, v := range userValidation {
		validate := validator.New()
		err := validate.Var(i, v)
		var (
			min         int
			max         int
			messageList []string
		)
		fmt.Sscanf(v, "required,min=%d,max=%d", &min, &max)
		fmt.Println(max)
		fmt.Println(min)
		if err != nil {
			if strings.Contains(err.Error(), "required") {
				message = fmt.Sprintf("%s %s", i, errorMapp["required"])
			}

			if strings.Contains(err.Error(), "max") {
				message = fmt.Sprintf("%s %s %v", i, errorMapp["max"], max)
			}

			if strings.Contains(err.Error(), "min") {
				message = fmt.Sprintf("%s %s %v", i, errorMapp["min"], min)
			}

			isValid = false
			fmt.Println(message)
			messageList = append(messageList, message)

		}

		errList[i] = append(messageList, messageList...)

		switch i {
		case "phone":
			isValid, errList = PhoneValidation(user.Phone, errList)
		case "password":
			isValid, errList = PasswordValidation(user.Password, errList)
		}

	}

	return isValid, errList
}

func PhoneValidation(phone string, errList map[string][]string) (bool, map[string][]string) {
	var messageList []string
	isValid := true
	message := ""
	pattern := regexp.MustCompile(`^\+62\d{8,11}$`)
	matched := pattern.MatchString(phone)

	if !matched {
		isValid = false
		message = fmt.Sprintf(`%s`, "must start with +62")
		messageList = append(messageList, message)
	}

	errList["phone"] = append(messageList, messageList...)
	return isValid, errList

}

func PasswordValidation(password string, errList map[string][]string) (bool, map[string][]string) {
	var messageList []string
	isValid := true
	message := ""
	pattern := regexp.MustCompile(`[A-Z]`)
	containsCapital := pattern.MatchString(password)

	if !containsCapital {
		isValid = false
		message = fmt.Sprintf(`%s`, "contains at least 1 Capital")
		messageList = append(messageList, message)
	}

	pattern = regexp.MustCompile(`\d`)
	containsNumber := pattern.MatchString(password)

	if !containsNumber {
		isValid = false
		message = "contains at least 1 number"
		messageList = append(messageList, message)
	}

	pattern = regexp.MustCompile(`\W`)
	containsSpecialChar := pattern.MatchString(password)
	if !containsSpecialChar {
		isValid = false
		message = "contains at least 1 spesial character"
		messageList = append(messageList, message)
	}

	errList["password"] = append(messageList, messageList...)
	return isValid, errList
}
