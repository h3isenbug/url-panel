package user

import "errors"

type User struct {
	Username string `json:"username" db:"username"`
	EMail    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
}

var (
	ErrUsernameTaken = errors.New("username is taken")
	ErrEMailTaken    = errors.New("email is taken")
)

type UserRepository interface {
	GetUserByUsername(username string) (*User, error)
	GetUserByEMail(email string) (*User, error)
	AddUser(username, email, password string) error
	//TODO UpdatePassword(password string) error
	//TODO ActivateUser(username string) error
}
