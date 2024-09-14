package auth

import (
	"code-hw/internal/storage"
	"context"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func BasicAuth(storage storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			headerParts := strings.SplitN(authHeader, " ", 2)
			if len(headerParts) != 2 || headerParts[0] != "Basic" {
				http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
				return
			}

			payload, _ := base64.StdEncoding.DecodeString(headerParts[1])
			pair := strings.SplitN(string(payload), ":", 2)
			if len(pair) != 2 {
				http.Error(w, "Неавторизованный пользователь", http.StatusUnauthorized)
				return
			}
			username, password := pair[0], pair[1]

			user, err := storage.GetUserByUsername(username)
			if err != nil {
				http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
				return
			}
			if user == nil || !checkPasswordHash(password, user.Password) {
				http.Error(w, "Неавторизованный пользователь", http.StatusUnauthorized)
				return
			}

			newCtx := context.WithValue(r.Context(), "userID", user.ID)

			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
