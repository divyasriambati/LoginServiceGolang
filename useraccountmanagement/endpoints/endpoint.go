package useraccountmanagement

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"../db"
	types "../models"
)

// **************************** Login ************************************

func handleLogin(w http.ResponseWriter, r *http.Request) {

	var user types.Login
	erro := json.NewDecoder(r.Body).Decode(&user)
	if erro != nil {
		http.Error(w, erro.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user exists in the database
	stmt, err := db.Prepare("SELECT username, password FROM `golangtestdb`.`userauth` WHERE username = ? AND password = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var dbUser types.Login
	err = stmt.QueryRow(user.Username, user.Password).Scan(&dbUser.Username, &dbUser.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Serialize the user data into a JSON response
	jsonUser, err := json.Marshal(dbUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonUser))
	w.Write([]byte("Login successfully"))

}

// ****************************sign up *************************************

func handleSignup(w http.ResponseWriter, r *http.Request) {
	var user types.SignupForm
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// fmt.Println(user)
	db, err := sql.Open("mysql", "root:Divya@123@tcp(localhost:3306)/golangtestdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO `golangtestdb`.`userauth` (fullname,username,email,password) VALUES (?, ?,?,?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	if user.Password == user.ConfirmPassword {
		_, err = stmt.Exec(user.FirstName+user.LastName, user.UserName, user.Email, user.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "The password and confirmation password do not match", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("sign up successfully"))
}

// ******************** DELETE USER**********************

func deleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var user types.DeleteUser
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Delete user from database
		result, err := db.Exec("DELETE FROM `golangtestdb`.`userauth` WHERE username=?", user.UserName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if user was deleted
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "User deleted successfully")
	}
}

// ******************************* update password *************************************

func updateUserPasswordHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var user types.UpdatePassword
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Verify old password
		var actualPassword string
		err = db.QueryRow("SELECT password FROM `golangtestdb`.`userauth` WHERE username=?", user.Username).Scan(&actualPassword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if actualPassword != user.OldPassword {
			http.Error(w, "Incorrect old password", http.StatusUnauthorized)
			return
		}

		// Update password in database
		_, err = db.Exec("UPDATE `golangtestdb`.`userauth` SET password=? WHERE username=?", user.NewPassword, user.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Password updated successfully")
	}
}

// ************************* Get Users *************************

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:Divya@123@tcp(localhost:3306)/golangtestdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT fullname, username, email, password FROM `golangtestdb`.`userauth`")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []types.GetUserDetails
	for rows.Next() {
		var user types.GetUserDetails
		err := rows.Scan(&user.FullName, &user.UserName, &user.Email, &user.Password)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
