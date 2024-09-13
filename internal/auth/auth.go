package auth

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type User struct {
	ID       int
	Username string
	Password string
}

type Storage interface {
	GetUserByUsername(username string) (*User, error)
}

// TODO реализовать авторизацию
func BasicAuth(storage Storage) func(http.Handler) http.Handler {
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

			// Проверяем логин и пароль
			user, err := storage.GetUserByUsername(username)
			if err != nil {
				http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
				return
			}
			if user == nil || !checkPasswordHash(password, user.Password) {
				http.Error(w, "Неавторизованный пользователь", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
