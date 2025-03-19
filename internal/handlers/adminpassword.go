package handlers

import (
	"log"
	"net/http"

	"github.com/1ssk/MicroHack.git/internal/database"
	"golang.org/x/crypto/bcrypt"
)

// ChangePasswordData представляет данные для формы смены пароля
type ChangePasswordData struct {
	ErrorMessage string
}

// ChangePasswordHandler обрабатывает запросы на смену пароля
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		// Отображаем форму смены пароля
		if err := renderTemplate(w, "change_password.html", nil); err != nil {
			handleError(w, err, http.StatusInternalServerError)
		}
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем данные из формы
	username := "admin" // Мы меняем пароль только для admin
	r.ParseForm()       // Добавляем разбор формы
	oldPassword := r.FormValue("old_password")
	newPassword := r.FormValue("new_password")
	confirmNewPassword := r.FormValue("confirm_new_password")

	// Валидация данных
	if oldPassword == "" || newPassword == "" || confirmNewPassword == "" {
		renderChangePasswordPage(w, "Все поля обязательны для заполнения")
		return
	}

	if newPassword != confirmNewPassword {
		renderChangePasswordPage(w, "Новые пароли не совпадают")
		return
	}

	// Получаем текущий хэшированный пароль admin из базы данных
	var hashedPassword string
	err := database.DB.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&hashedPassword)
	if err != nil {
		log.Printf("Ошибка при получении пароля пользователя: %v", err)
		renderChangePasswordPage(w, "Ошибка при смене пароля")
		return
	}

	// Проверяем, совпадает ли старый пароль
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(oldPassword))
	if err != nil {
		renderChangePasswordPage(w, "Неверный старый пароль")
		return
	}

	// Хэшируем новый пароль
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка при хэшировании нового пароля: %v", err)
		renderChangePasswordPage(w, "Ошибка при смене пароля")
		return
	}

	// Обновляем пароль в базе данных
	_, err = database.DB.Exec("UPDATE users SET password = $1 WHERE username = $2", newHashedPassword, username)
	if err != nil {
		log.Printf("Ошибка при обновлении пароля: %v", err)
		renderChangePasswordPage(w, "Ошибка при смене пароля")
		return
	}

	// Успешная смена пароля
	http.Redirect(w, r, "/admin", http.StatusFound)
}

// renderChangePasswordPage рендерит страницу смены пароля с ошибкой
func renderChangePasswordPage(w http.ResponseWriter, errorMessage string) {
	data := ChangePasswordData{
		ErrorMessage: errorMessage,
	}
	if err := renderTemplate(w, "change_password.html", data); err != nil {
		handleError(w, err, http.StatusInternalServerError)
	}
}
