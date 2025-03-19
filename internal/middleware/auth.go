package middleware

import (
	"net/http"
	"strings"

	"github.com/1ssk/MicroHack.git/internal/handlers"
	"github.com/dgrijalva/jwt-go"
)

// AuthMiddleware проверяет JWT-токен и роль пользователя
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем куку с токеном
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Декодируем токен
		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return handlers.JWTKey, nil
		})

		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Проверяем доступ к /register
		if strings.HasPrefix(r.URL.Path, "/register") && claims.Role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Проверяем доступ к /admin
		if strings.HasPrefix(r.URL.Path, "/admin") && claims.Role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Передаем управление дальше
		next.ServeHTTP(w, r)
	})
}
