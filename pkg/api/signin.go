package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("TODO_PASSWORD"))

// auth используется для проверки аутентификации
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if len(jwtSecret) == 0 {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			writeJson(w, map[string]string{"error": "Требуется аутентификация"}, http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			writeJson(w, map[string]string{"error": "Недействительный токен"}, http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// signinHandler отправляет POST-запрос для входа в систему
func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, map[string]string{"error": "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJson(w, map[string]string{"error": "Неверный формат запроса"}, http.StatusBadRequest)
		return
	}

	if request.Password != string(jwtSecret) {
		writeJson(w, map[string]string{"error": "Неверный пароль"}, http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash":    string(jwtSecret),
		"expires": time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		writeJson(w, map[string]string{"error": "Ошибка генерации токена"}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{"token": tokenString}, http.StatusOK)
}
