package auth

import (
	"errors"
	"github.com/h3isenbug/url-panel/repositories/jwt"
	userRepository "github.com/h3isenbug/url-panel/repositories/user"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService interface {
	LoginWithUsername(username, password string) (string, time.Time, error)
	LoginWithEMail(email, password string) (string, time.Time, error)
	Register(username, email, password string) error
	ParseToken(token string) (*jwt.Token, error)
	//TODO Logout(authToken, refreshToken string) error
	//TODO VerifyEMail(token string) error
	//TODO ChangePassword(username, currentPassword, newPassword string) error
	//TODO ForgotPassword(username ) error
}

var (
	ErrWrongCredentials = errors.New("wrong credentials")
)

type AuthServiceV1 struct {
	userRepo userRepository.UserRepository
	loginJWT jwt.JWTWrapper
}

func NewAuthServiceV1(userRepo userRepository.UserRepository, jwtKey []byte) AuthService {
	return &AuthServiceV1{userRepo: userRepo, loginJWT: jwt.NewJWTGoWrapper(time.Minute*20, "login", jwtKey)}
}

func (service AuthServiceV1) LoginWithUsername(username, password string) (string, time.Time, error) {
	user, err := service.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", time.Time{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", time.Time{}, err
	}

	token, expiresAt, err := service.loginJWT.CreateToken(user.EMail)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func (service AuthServiceV1) LoginWithEMail(email, password string) (string, time.Time, error) {
	user, err := service.userRepo.GetUserByEMail(email)
	if err != nil {
		return "", time.Time{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", time.Time{}, nil
	}

	token, expiresAt, err := service.loginJWT.CreateToken(email)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func (service AuthServiceV1) Register(username, email, password string) error {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := service.userRepo.AddUser(username, email, string(hashBytes)); err != nil {
		return err
	}

	//TODO tell workers to send verification email
	return nil
}

func (service AuthServiceV1) ParseToken(token string) (*jwt.Token, error) {
	return service.loginJWT.ParseToken(token)
}
