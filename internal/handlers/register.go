package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/1ssk/MicroHack.git/internal/database"
	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler обрабатывает запросы на страницу регистрации
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Отображаем форму регистрации
		if err := renderTemplate(w, "register.html", nil); err != nil {
			handleError(w, err, http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		// Получаем данные из формы
		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm-password")

		// Валидация данных
		if username == "" || password == "" || confirmPassword == "" {
			renderRegisterPage(w, "Все поля обязательны для заполнения")
			return
		}

		if password != confirmPassword {
			renderRegisterPage(w, "Пароли не совпадают")
			return
		}

		// Проверяем, не существует ли уже пользователь с таким именем
		var count int
		err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
		if err != nil {
			log.Printf("Ошибка при проверке существования пользователя: %v", err)
			renderRegisterPage(w, "Ошибка при регистрации")
			return
		}
		if count > 0 {
			renderRegisterPage(w, "Пользователь с таким именем уже существует")
			return
		}

		// Хэшируем пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Ошибка при хешировании пароля: %v", err)
			renderRegisterPage(w, "Ошибка при регистрации")
			return
		}

		// Сохраняем пользователя в базу данных
		_, err = database.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, hashedPassword)
		if err != nil {
			log.Printf("Ошибка при сохранении пользователя: %v", err)
			renderRegisterPage(w, "Ошибка при регистрации")
			return
		}

		// Успешная регистрация — перенаправляем на страницу логина
		fmt.Println("register seccesful")
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// renderRegisterPage рендерит страницу регистрации с ошибкой
func renderRegisterPage(w http.ResponseWriter, errorMessage string) {
	data := map[string]interface{}{
		"ErrorMessage": errorMessage,
	}
	if err := renderTemplate(w, "register.html", data); err != nil {
		handleError(w, err, http.StatusInternalServerError)
	}
}
