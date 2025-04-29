package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

type SignInRequest struct {
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if req.Password != pass {
		writeJSON(w, http.StatusUnauthorized, SignInResponse{Error: "Неверный пароль"})
		return
	}

	token, err := generateJWT(req.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, SignInResponse{Error: "Ошибка создания токена"})
		return
	}

	writeJSON(w, http.StatusOK, SignInResponse{Token: token})
	fmt.Println("Token:", token)
}

func generateJWT(password string) (string, error) {
	claims := jwt.MapClaims{
		"passwordHash": fmt.Sprintf("%x", password),
		"exp":          jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
