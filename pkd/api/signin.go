package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type signinRes struct {
	Token string `json:"token"`
}
type signinReq struct {
	Password string `json:"password"`
}

var secretKey = []byte("secret_key")

var pass string

func InitPassword(password string) {
	pass = password
}

func checksumPassword(pass string) string {
	hash := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(hash[:])
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(pass) > 0 {
			var jwtCookie string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtCookie = cookie.Value
			}
			token, err := jwt.Parse(jwtCookie, func(t *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})
			if err != nil {
				http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
				return
			}
			if !token.Valid {
				http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
				return
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				receivedChecksum := claims["checksum"].(string)
				currentChecksum := checksumPassword(pass)
				if receivedChecksum != currentChecksum {
					http.Error(w, "Пароль изменился и токен недействителен", http.StatusUnauthorized)
					return
				}
			} else {
				http.Error(w, "Неверная структура токена ", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
func signinHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respErrJSON(w, "Ошибка чтения тела", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var signinReq signinReq
	if err = json.Unmarshal(body, &signinReq); err != nil {
		respErrJSON(w, "Ошибка десириализации", http.StatusBadRequest)
		return
	}
	if signinReq.Password != pass {
		respErrJSON(w, "Пароль неверный", http.StatusBadRequest)
		return
	}
	checksum := checksumPassword(pass)
	claims := jwt.MapClaims{
		"exp":      time.Now().Add(8 * time.Hour).Unix(),
		"checksum": checksum,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString(secretKey)
	if err != nil {
		respErrJSON(w, "Не удалось подписать токен", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(signinRes{Token: signedToken}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
