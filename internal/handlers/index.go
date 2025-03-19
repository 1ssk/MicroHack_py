package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/1ssk/MicroHack.git/internal/database"
	"github.com/golang-jwt/jwt"
)

// IndexHandler обрабатывает запросы на главную страницу
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Отображаем главную страницу
	if err := renderTemplate(w, "index.html", nil); err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

// Course - структура курса
type Course struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// MyCoursesHandler обрабатывает запросы к API /api/my-courses
func MyCoursesHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем куку с токеном
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Println("Ошибка: нет куки с токеном")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Декодируем JWT-токен
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return JWTKey, nil
	})
	if err != nil || !token.Valid {
		log.Println("Ошибка: токен недействителен")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем user_id по username
	var userID int
	err = database.DB.QueryRow("SELECT id FROM users WHERE username = $1", claims.Username).Scan(&userID)
	if err != nil {
		log.Printf("Ошибка поиска user_id для %s: %v", claims.Username, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Получаем курсы пользователя из базы через промежуточную таблицу
	rows, err := database.DB.Query(`
		SELECT c.id, c.title
		FROM courses c
		JOIN user_courses uc ON c.id = uc.course_id
		WHERE uc.user_id = $1
	`, userID)
	if err != nil {
		log.Printf("Ошибка запроса курсов для user_id %d: %v", userID, err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.ID, &course.Title); err != nil {
			log.Printf("Ошибка при сканировании данных курсов: %v", err)
			http.Error(w, "Database scan error", http.StatusInternalServerError)
			return
		}
		courses = append(courses, course)
	}

	// Проверяем ошибки после цикла
	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при итерации по курсам: %v", err)
		http.Error(w, "Database iteration error", http.StatusInternalServerError)
		return
	}

	// Отправляем курсы пользователю
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}
