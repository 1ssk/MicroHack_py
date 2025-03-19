package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/1ssk/MicroHack.git/internal/database"
	"golang.org/x/crypto/bcrypt"
)

// LoginHandler обрабатывает запросы на страницу логина
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if err := renderTemplate(w, "login.html", nil); err != nil {
			handleError(w, err, http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var storedPassword string
		var role string // Добавляем переменную для хранения роли
		err := database.DB.QueryRow("SELECT password, role FROM users WHERE username = $1", username).Scan(&storedPassword, &role)
		if err != nil {
			renderLoginPage(w, "Пользователь не найден")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
		if err == nil {
			// Генерация JWT-токена
			token, err := generateJWTToken(username, role)
			if err != nil {
				handleError(w, err, http.StatusInternalServerError)
				return
			}

			// Устанавливаем токен в куки
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   token,
				Expires: time.Now().Add(24 * time.Hour), // Время жизни куки
			})

			// Перенаправляем на защищенную страницу в зависимости от роли
			switch role {
			case "admin":
				http.Redirect(w, r, "/admin", http.StatusFound)
			case "user":
				http.Redirect(w, r, "/", http.StatusFound)
			default:
				log.Printf("Неизвестная роль пользователя: %s", role)
				http.Redirect(w, r, "/", http.StatusFound)
			}

		} else {
			renderLoginPage(w, "Неверный пароль")
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// renderLoginPage рендерит страницу логина с ошибкой
func renderLoginPage(w http.ResponseWriter, errorMessage string) {
	data := map[string]interface{}{
		"ErrorMessage": errorMessage,
	}
	if err := renderTemplate(w, "login.html", data); err != nil {
		handleError(w, err, http.StatusInternalServerError)
	}
}
