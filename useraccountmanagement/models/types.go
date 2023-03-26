package types

type SignupForm struct {
	FirstName       string `json:"firstname"`
	LastName        string `json:"lastname"`
	UserName        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmpassword"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DeleteUser struct {
	UserName string `json:"username"`
}

type UpdatePassword struct {
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type GetUserDetails struct {
	FullName string `json:"fullname"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
