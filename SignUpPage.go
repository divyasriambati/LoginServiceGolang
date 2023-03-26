package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
	Fullname        string `json:"fullname"`
	Contactno       string `json:"contactno"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	PWord           string `json:"password"`
	OldPassword     string `json:"oldpassword"`
	NewPassword     string `json: "newpassword"`
	ConfirmPassword string `json:"confirmpassword"`
}

func main() {
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", "root:Chandu@9@tcp(localhost:3306)/users")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Connected!")
	http.HandleFunc("/", handler)
	http.HandleFunc("/SignUp", SignUpRequest)
	http.HandleFunc("/Login", LoginRequest)
	http.HandleFunc("/Update", UpdateRequest)
	http.HandleFunc("/Delete", DeleteRequest)
	log.Printf("Starting the server at the port 9999")
	err = http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Println("Error Occured", err)
	}
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Successsssssssss")
}

// SignUp Request
func SignUpRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
	//Get the Data
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Decoding Error", http.StatusInternalServerError)
		return
	}
	fmt.Println(user.Fullname)
	//Validate the data
	if user.Fullname == " " || user.Fullname == "" {
		http.Error(w, "FullName is Required", http.StatusBadRequest)
		return
	}
	if user.Username == " " || user.Username == "" {
		http.Error(w, "UserName is Required", http.StatusBadRequest)
		return
	}
	if user.PWord == " " || user.PWord == "" {
		http.Error(w, "Password is Required", http.StatusBadRequest)
		return
	}
	if len(user.PWord) < 8 {
		http.Error(w, "Password is Not Valid", http.StatusBadRequest)
		return
	}
	if user.ConfirmPassword != user.PWord {
		http.Error(w, "Passwords don't match", http.StatusBadRequest)
		return
	}
	//Execute the query
	query, err := db.Prepare(`INSERT INTO users VALUES(?,?,?,?,?)`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = query.Exec(user.Fullname, user.Contactno, user.Email, user.Username, user.PWord)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Successfully Created new user: %s", user.Username)
	w.WriteHeader(http.StatusCreated)

}

// Login with UserName and Password
func LoginRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
	//Get the Data
	var loginData User
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Login Decoding Error", http.StatusInternalServerError)
		return
	}
	fmt.Printf(loginData.Username)
	//Validate the Data
	query, err := db.Query(`select count(*), PWord from users.users where UserName = ? Group By UserName`, loginData.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var count int
	var passw string
	for query.Next() {
		if err := query.Scan(&count, &passw); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if count == 0 {
		http.Error(w, "Username Not Found", http.StatusBadRequest)
		return
	}
	if loginData.PWord != passw {
		http.Error(w, "Incorrect Password", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Login Successful")
	w.WriteHeader(http.StatusCreated)
}

// Update the Password
func UpdateRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
	//Get the Data
	var updateUser User
	err := json.NewDecoder(r.Body).Decode(&updateUser)
	if err != nil {
		http.Error(w, "Updating Error", http.StatusInternalServerError)
		return
	}
	//Validate the Data
	query, err := db.Query(`select count(*), PWord from users.users where UserName = ? Group By UserName`, updateUser.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var count int
	var passw string
	for query.Next() {
		if err := query.Scan(&count, &passw); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if count == 0 {
		http.Error(w, "Username Not Found", http.StatusBadRequest)
		return
	}
	if passw != updateUser.OldPassword {
		http.Error(w, "Incorrect old Password", http.StatusBadRequest)
		fmt.Println(updateUser.OldPassword, updateUser.PWord)
		return
	}
	//fmt.Println(updateUser.OldPassword, updateUser.PWord)
	if len(updateUser.NewPassword) < 8 {
		http.Error(w, "Password is Not Valid", http.StatusBadRequest)
		return
	}
	if updateUser.ConfirmPassword != updateUser.NewPassword {
		http.Error(w, "Passwords don't match", http.StatusBadRequest)
		return
	}
	//Update the Data
	_, err = db.Query(`Update users.users set PWord = ? where UserName = ? `, updateUser.NewPassword, updateUser.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Updated the password for user: %s", updateUser.Username)
	w.WriteHeader(http.StatusAccepted)
}

// Delete the Data
func DeleteRequest(w http.ResponseWriter, r *http.Request) {
	//Validate the Method
	if r.Method != "DELETE" {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
	//Get the Query Params
	username := r.URL.Query().Get("username")
	query, err := db.Query(`select count(*), UserName from users.users where UserName = ? Group By UserName`, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Execute the Query
	var count int
	var uname string
	for query.Next() {
		if err := query.Scan(&count, &uname); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if count == 0 {
		//fmt.Println(count, username, uname)
		http.Error(w, "User not Found", http.StatusBadRequest)
		return
	}
	//Delete the data
	_, err = db.Query(`DELETE FROM users.users where UserName = ?`, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Deleted the username from DataBase")
	w.WriteHeader(http.StatusAccepted)
}
