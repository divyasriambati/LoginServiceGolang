package useraccountmanagement

import (
	"encoding/json"
	"fmt"
	"net/http"

	userdb "github.com/divyasriambati/LoginServiceGolang/useraccountmanagement/db"
	types "github.com/divyasriambati/LoginServiceGolang/useraccountmanagement/models"
	validators "github.com/divyasriambati/LoginServiceGolang/useraccountmanagement/validations"
)

// **************************** Login ************************************

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	var user types.Login

	// extract user details from the request
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user exists in the database

	err = userdb.GetUserLoginDetails(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Login successfully"))

}

// ****************************sign up *************************************

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	var user types.User

	// extract user details from the request
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate user details
	status, err := validators.ValidateUser(user)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	// insert user details into db
	err = userdb.InsertUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("sign up successfully"))
}

// ******************** DELETE USER**********************

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var user types.DeleteUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//delete user from DB
	err = userdb.DeleteUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User deleted successfully")
}

// ******************************* update password *************************************

func UpdateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var user types.UpdatePassword
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate update password details
	status, err := validators.ValidateUpdatePassword(user)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	// update user password in db
	err = userdb.UpdateUserPassword(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Password updated successfully")
}

// ************************* Get Users *************************

func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	var user types.GetUserDetails

	//get the user details
	users, err := userdb.GetUsers(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
