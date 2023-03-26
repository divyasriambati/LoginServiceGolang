package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	useraccountmanagement "github.com/divyasriambati/LoginServiceGolang/useraccountmanagement/endpoints"
)

func main() {

	http.HandleFunc("/signup", useraccountmanagement.HandleSignup)
	http.HandleFunc("/login", useraccountmanagement.HandleLogin)
	http.HandleFunc("/deleteUser", useraccountmanagement.DeleteUserHandler)
	http.HandleFunc("/update", useraccountmanagement.UpdateUserPasswordHandler)
	http.HandleFunc("/getusers", useraccountmanagement.HandleGetUsers)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
