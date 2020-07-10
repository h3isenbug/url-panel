package auth

import (
	"encoding/json"
	"errors"
	log2 "github.com/h3isenbug/url-common/pkg/services/log"
	"github.com/h3isenbug/url-panel/handlers"
	"github.com/h3isenbug/url-panel/repositories"
	"github.com/h3isenbug/url-panel/repositories/user"
	"github.com/h3isenbug/url-panel/services/auth"
	"net/http"
	"strings"
	"time"
)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	//TODO ActivateAccount
	//TODO ChangePassword
	//TODO Logout
}

type AuthHandlerV1 struct {
	authService auth.AuthService
	logService  log2.LogService
}

func NewAuthHandlerV1(authService auth.AuthService, logService log2.LogService) *AuthHandlerV1 {
	return &AuthHandlerV1{authService: authService, logService: logService}
}

func (handler AuthHandlerV1) Login(w http.ResponseWriter, r *http.Request) {
	var postFields struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&postFields); err != nil {
		handlers.SendError(w, http.StatusBadRequest)
		return
	}

	var err error
	var token string
	var expiresAt time.Time
	if strings.Contains(postFields.Identifier, "@") {
		token, expiresAt, err = handler.authService.LoginWithEMail(postFields.Identifier, postFields.Password)
	} else {
		token, expiresAt, err = handler.authService.LoginWithUsername(postFields.Identifier, postFields.Password)
	}
	if errors.Is(err, auth.ErrWrongCredentials) || errors.Is(err, repositories.ErrNotFound) {
		handlers.SendErrorWithCustomMessage(w, http.StatusBadRequest, "wrong credentials")
		return
	}
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError)
		handler.logService.Error("could not authenticate user(%s): %s", postFields.Identifier, err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, struct {
		Token     string `json:"token"`
		ExpiresAt string `json:"expiresAt"`
	}{
		Token:     token,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	})
}

func (handler AuthHandlerV1) Register(w http.ResponseWriter, r *http.Request) {
	var postFields struct {
		Username string `json:"username"`
		EMail    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&postFields); err != nil {
		handlers.SendError(w, http.StatusBadRequest)
		return
	}

	err := handler.authService.Register(postFields.Username, postFields.EMail, postFields.Password)
	if errors.Is(err, user.ErrUsernameTaken) {
		handlers.SendErrorWithCustomMessage(w, http.StatusBadRequest, "that username is taken")
		return
	}
	if errors.Is(err, user.ErrEMailTaken) {
		handlers.SendErrorWithCustomMessage(w, http.StatusBadRequest, "a user is registered with that email")
		return
	}
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError)
		return
	}

	handlers.SendResponse(w, http.StatusOK, nil)
}
