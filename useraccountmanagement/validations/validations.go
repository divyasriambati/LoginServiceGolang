package validators

import (
	"errors"
	"net/http"
	"regexp"

	types "github.com/divyasriambati/LoginServiceGolang/useraccountmanagement/models"
)

func IsEmailValid(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

func ValidateUser(user types.User) (int, error) {

	if user.FirstName == "" {
		return http.StatusBadRequest, errors.New("Firstname missing")
	}
	if user.LastName == "" {
		return http.StatusBadRequest, errors.New("Lastname missing")
	}
	if user.UserName == "" {
		return http.StatusBadRequest, errors.New("Username missing")
	}
	if IsEmailValid(user.Email) {
		return http.StatusBadRequest, errors.New("Invalid Email")
	}
	if len(user.Password) < 8 {
		return http.StatusBadRequest, errors.New("Password missing")

	}
	if len(user.ConfirmPassword) < 8 {
		return http.StatusBadRequest, errors.New("Confirm password missing")

	}
	if user.Password != user.ConfirmPassword {
		return http.StatusBadRequest, errors.New("passwords donot match")

	}
	return 200, nil
}

func ValidateUpdatePassword(user types.UpdatePassword) (int, error) {

	if user.Username == "" {
		return http.StatusBadRequest, errors.New("Username missing")
	}
	if len(user.OldPassword) < 8 {
		return http.StatusBadRequest, errors.New("Old Password missing")

	}
	if len(user.NewPassword) < 8 {
		return http.StatusBadRequest, errors.New("New password missing")
	}

	return 200, nil
}
