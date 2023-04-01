package main

import (
	"log"
	"net/http"

	useraccountmanagement "github.com/divyasriambati/LoginServiceGolang/useraccountmanagement/endpoints"
	auth "github.com/divyasriambati/LoginServiceGolang/useraccountmanagement/validations"
)

func main() {

	http.HandleFunc("/signup", useraccountmanagement.HandleSignup)
	http.HandleFunc("/login", useraccountmanagement.HandleLogin)
	http.HandleFunc("/deleteUser", auth.Authenticate(useraccountmanagement.DeleteUserHandler))
	http.HandleFunc("/update", auth.Authenticate(useraccountmanagement.UpdateUserPasswordHandler))
	http.HandleFunc("/getusers", useraccountmanagement.HandleGetUsers)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
