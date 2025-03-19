package handlers

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims — структура для хранения данных в JWT-токене
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"` // Добавлено поле роли
	jwt.StandardClaims
}

// Секретный ключ для подписи JWT-токена
var JWTKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// generateJWTToken создает JWT-токен для указанного пользователя и его роли
func generateJWTToken(username, role string) (string, error) {
	// Время жизни токена (например, 24 часа)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Создаем claims (данные для токена)
	claims := &Claims{
		Username: username,
		Role:     role, // Записываем роль пользователя в токен
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // Время истечения токена
		},
	}

	// Создаем токен с алгоритмом подписи HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен с использованием секретного ключа
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
