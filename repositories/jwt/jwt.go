package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

type Token struct {
	ExpiresAt time.Time
	ID        string
	IssuedAt  time.Time
	Issuer    string
	NotBefore time.Time
	Subject   string
}

type JWTWrapper interface {
	CreateToken(username string) (string, time.Time, error)
	ParseToken(token string) (*Token, error)
}

type JWTGoWrapper struct {
	tokenLifespan time.Duration
	issuer        string
	key           []byte
}

func NewJWTGoWrapper(tokenLifespan time.Duration, issuer string, key []byte) JWTWrapper {
	return &JWTGoWrapper{tokenLifespan: tokenLifespan, issuer: issuer, key: key}
}

func (wrapper JWTGoWrapper) CreateToken(email string) (string, time.Time, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", time.Time{}, err
	}
	var expiresAt = time.Now().Add(wrapper.tokenLifespan)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		ExpiresAt: expiresAt.Unix(),
		Id:        id.String(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    wrapper.issuer,
		NotBefore: time.Now().Unix(),
		Subject:   email,
	}).SignedString(wrapper.key)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func (wrapper JWTGoWrapper) ParseToken(tokenString string) (*Token, error) {
	var claims jwt.StandardClaims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}

		return wrapper.key, nil
	})
	if err != nil {
		return nil, err
	}

	if claims.Valid() != nil {
		return nil, errors.New("token is not valid now")
	}

	if !claims.VerifyIssuer(wrapper.issuer, true) {
		return nil, fmt.Errorf("token is issued by %s, not %s", claims.Issuer, wrapper.issuer)
	}

	return &Token{
		ExpiresAt: time.Unix(claims.ExpiresAt, 0),
		ID:        claims.Id,
		IssuedAt:  time.Unix(claims.IssuedAt, 0),
		Issuer:    claims.Issuer,
		NotBefore: time.Unix(claims.NotBefore, 0),
		Subject:   claims.Subject,
	}, nil
}
