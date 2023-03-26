package userdb

import (
	"errors"

	types "github.com/divyasriambati/LoginServiceGolang/useraccountmanagement/models"
)

func InsertUser(user types.User) error {

	db, err := Connectdb()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO `golangtestdb`.`userauth` (fullname,username,email,password) VALUES (?, ?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName+user.LastName, user.UserName, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func GetUserLoginDetails(user types.Login) error {

	db, err := Connectdb()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT username, password FROM `golangtestdb`.`userauth` WHERE username = ? AND password = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var dbUser types.Login
	err = stmt.QueryRow(user.Username, user.Password).Scan(&dbUser.Username, &dbUser.Password)
	if err != nil {
		// return errors.New("Invalid Email or password")
		return err
	}

	return nil
}

func DeleteUser(user types.DeleteUser) error {

	db, err := Connectdb()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Delete user from database
	result, err := db.Exec("DELETE FROM `golangtestdb`.`userauth` WHERE username=?", user.UserName)
	if err != nil {
		return err
	}

	// Check if user was deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("User not found")
	}

	return nil
}

func UpdateUserPassword(user types.UpdatePassword) error {

	db, err := Connectdb()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Verify old password
	var actualPassword string
	err = db.QueryRow("SELECT password FROM `golangtestdb`.`userauth` WHERE username=?", user.Username).Scan(&actualPassword)
	if err != nil {
		return err
	}
	if actualPassword != user.OldPassword {
		// http.Error(w, "Incorrect old password", http.StatusUnauthorized)
		return errors.New("Incorrect old password")
	}

	// Update password in database
	_, err = db.Exec("UPDATE `golangtestdb`.`userauth` SET password=? WHERE username=?", user.NewPassword, user.Username)
	if err != nil {
		return err
	}

	return nil
}

func GetUsers(user types.GetUserDetails) ([]types.GetUserDetails, error) {
	db, err := Connectdb()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT fullname, username, email, password FROM `golangtestdb`.`userauth`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.GetUserDetails
	for rows.Next() {
		var user types.GetUserDetails
		err := rows.Scan(&user.FullName, &user.UserName, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
