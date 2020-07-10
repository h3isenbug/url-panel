package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/h3isenbug/url-panel/services/auth"
	"net/http"
)

const (
	ContextKeyURLParams = iota
	ContextKeyAuthToken
	ContextKeyEMail
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func GorillaMuxURLParamMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ContextKeyURLParams, mux.Vars(r))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthTokenMiddleware(authService auth.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("AuthToken")
			if len(tokenString) == 0 {
				SendError(w, http.StatusForbidden)
				return
			}
			token, err := authService.ParseToken(tokenString)
			if err != nil {
				SendError(w, http.StatusForbidden)
				return
			}
			ctxWithEMail := context.WithValue(r.Context(), ContextKeyEMail, token.Subject)
			next.ServeHTTP(w, r.WithContext(ctxWithEMail))
		})
	}
}

func GetURLParams(r *http.Request) map[string]string {
	value := r.Context().Value(ContextKeyURLParams)
	if value == nil {
		return nil
	}

	return value.(map[string]string)
}

func SendError(w http.ResponseWriter, statusCode int) error {
	return SendErrorWithCustomMessage(w, statusCode, http.StatusText(statusCode))
}

func SendErrorWithCustomMessage(w http.ResponseWriter, statusCode int, message string) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(&Response{
		Message: message,
	})
}

func SendResponseWithMessage(w http.ResponseWriter, statusCode int, message string, result interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(&Response{
		Data:    result,
		Message: message,
	})
}

func SendResponse(w http.ResponseWriter, statusCode int, result interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(&Response{
		Data: result,
	})
}
