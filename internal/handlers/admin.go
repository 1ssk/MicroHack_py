package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/1ssk/MicroHack.git/internal/database"
	"github.com/gorilla/mux"
)

// AdminHandler обрабатывает запросы на админ-панель
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем список пользователей
	users, err := GetUsers()
	if err != nil {
		http.Error(w, "Ошибка при получении пользователей", http.StatusInternalServerError)
		log.Println("Ошибка при получении пользователей:", err)
		return
	}

	// Получаем список курсов
	courses, err := GetCourses()
	if err != nil {
		http.Error(w, "Ошибка при получении курсов", http.StatusInternalServerError)
		log.Println("Ошибка при получении курсов:", err)
		return
	}

	// Получаем список уроков
	lessons, err := GetLessons()
	if err != nil {
		http.Error(w, "Ошибка при получении уроков", http.StatusInternalServerError)
		log.Println("Ошибка при получении уроков:", err)
		return
	}

	// Отображаем шаблон админ-панели с данными
	data := struct {
		Users   []User
		Courses []Course
		Lessons []Lesson
	}{
		Users:   users,
		Courses: courses,
		Lessons: lessons,
	}

	if err := renderTemplate(w, "admin.html", data); err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
}

// getUsers получает список пользователей
func GetUsers() ([]User, error) {
	rows, err := database.DB.Query("SELECT id, username FROM users")
	if err != nil {
		log.Println("Ошибка при выполнении запроса:", err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			log.Println("Ошибка при чтении строки результата:", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// getCourses получает список курсов
func GetCourses() ([]Course, error) {
	rows, err := database.DB.Query("SELECT id, title FROM courses")
	if err != nil {
		log.Println("Ошибка при выполнении запроса:", err)
		return nil, err
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		err := rows.Scan(&course.ID, &course.Title)
		if err != nil {
			log.Println("Ошибка при чтении строки результата:", err)
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}

// createCourse создает новый курс
func CreateCourse(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.FormValue("title")

	// Вставка нового курса в базу данных
	_, err := database.DB.Exec("INSERT INTO courses (title) VALUES ($1)", title)
	if err != nil {
		http.Error(w, "Ошибка при создании курса", http.StatusInternalServerError)
		log.Println("Ошибка при добавлении курса:", err)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// assignCourse добавляет курс пользователю
func AssignCourse(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userID := r.FormValue("user_id")
	courseID := r.FormValue("course_id")

	// Добавление связи между пользователем и курсом
	_, err := database.DB.Exec("INSERT INTO user_courses (user_id, course_id) VALUES ($1, $2)", userID, courseID)
	if err != nil {
		http.Error(w, "Ошибка при добавлении курса пользователю", http.StatusInternalServerError)
		log.Println("Ошибка при добавлении курса пользователю:", err)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// createLesson создает новый урок
func CreateLesson(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.FormValue("title")
	url := r.FormValue("url")
	courseID := r.FormValue("course_id")

	// Вставка нового урока в базу данных
	_, err := database.DB.Exec("INSERT INTO lessons (title, course_id, url) VALUES ($1, $2, $3)", title, courseID, url)
	if err != nil {
		http.Error(w, "Ошибка при создании урока", http.StatusInternalServerError)
		log.Println("Ошибка при добавлении урока:", err)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// GetUserCourses получает курсы, связанные с пользователем
func GetUserCourses(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	courses, err := GetCoursesByUser(userID)
	if err != nil {
		http.Error(w, "Ошибка при получении курсов пользователя", http.StatusInternalServerError)
		log.Println("Ошибка при получении курсов пользователя:", err)
		return
	}

	// Отправляем список курсов в формате JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(courses)
	if err != nil {
		http.Error(w, "Ошибка при отправке данных", http.StatusInternalServerError)
	}
}

// getCoursesByUser получает курсы, связанные с пользователем
func GetCoursesByUser(userID string) ([]Course, error) {
	rows, err := database.DB.Query("SELECT courses.id, courses.title FROM courses INNER JOIN user_courses ON courses.id = user_courses.course_id WHERE user_courses.user_id = $1", userID)
	if err != nil {
		log.Println("Ошибка при выполнении запроса:", err)
		return nil, err
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		err := rows.Scan(&course.ID, &course.Title)
		if err != nil {
			log.Println("Ошибка при чтении строки результата:", err)
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}

// User представляет пользователя
type User struct {
	ID       int
	Username string
}

func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	courseID := r.FormValue("course_id")
	if courseID == "" {
		http.Error(w, "Course ID is required", http.StatusBadRequest)
		return
	}

	// Начинаем транзакцию для безопасного удаления
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Ошибка при начале транзакции", http.StatusInternalServerError)
		log.Println("Ошибка при начале транзакции:", err)
		return
	}

	// Удаляем связи пользователей с курсом
	_, err = tx.Exec("DELETE FROM user_courses WHERE course_id = $1", courseID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Ошибка при удалении связей пользователей с курсом", http.StatusInternalServerError)
		log.Println("Ошибка при удалении связей пользователей с курсом:", err)
		return
	}

	// Удаляем уроки, связанные с курсом
	_, err = tx.Exec("DELETE FROM lessons WHERE course_id = $1", courseID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Ошибка при удалении уроков", http.StatusInternalServerError)
		log.Println("Ошибка при удалении уроков:", err)
		return
	}

	// Удаляем сам курс
	_, err = tx.Exec("DELETE FROM courses WHERE id = $1", courseID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Ошибка при удалении курса", http.StatusInternalServerError)
		log.Println("Ошибка при удалении курса:", err)
		return
	}

	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		http.Error(w, "Ошибка при фиксации транзакции", http.StatusInternalServerError)
		log.Println("Ошибка при фиксации транзакции:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteUser удаляет пользователя и связанные данные
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	userID := r.FormValue("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Начинаем транзакцию для безопасного удаления
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Ошибка при начале транзакции", http.StatusInternalServerError)
		log.Println("Ошибка при начале транзакции:", err)
		return
	}

	// Удаляем связи пользователя с курсами
	_, err = tx.Exec("DELETE FROM user_courses WHERE user_id = $1", userID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Ошибка при удалении связей пользователя с курсами", http.StatusInternalServerError)
		log.Println("Ошибка при удалении связей пользователя с курсами:", err)
		return
	}

	// Удаляем пользователя
	_, err = tx.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Ошибка при удалении пользователя", http.StatusInternalServerError)
		log.Println("Ошибка при удалении пользователя:", err)
		return
	}

	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		http.Error(w, "Ошибка при фиксации транзакции", http.StatusInternalServerError)
		log.Println("Ошибка при фиксации транзакции:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteLesson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	lessonID := r.FormValue("lesson_id")
	if lessonID == "" {
		http.Error(w, "Lesson ID is required", http.StatusBadRequest)
		return
	}

	// Удаляем урок
	_, err := database.DB.Exec("DELETE FROM lessons WHERE id = $1", lessonID)
	if err != nil {
		http.Error(w, "Ошибка при удалении урока", http.StatusInternalServerError)
		log.Println("Ошибка при удалении урока:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
