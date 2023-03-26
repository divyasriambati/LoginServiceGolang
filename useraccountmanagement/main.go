package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	http.HandleFunc("/signup", handleSignup)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/deleteUser", deleteUserHandler(db))
	http.HandleFunc("/update", updateUserPasswordHandler(db))
	http.HandleFunc("/getusers", handleGetUsers)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
