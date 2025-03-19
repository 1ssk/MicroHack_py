package handlers

import (
	"log"
	"net/http"

	"github.com/1ssk/MicroHack.git/internal/database"
	"github.com/gorilla/mux"
)

type Lesson struct {
	ID          int
	Title       string
	CourseID    int
	CourseTitle string
	URL         string
}

// LessonHandler обрабатывает запросы для отображения уроков по курсу
func LessonHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем course_id из параметров URL
	vars := mux.Vars(r)
	courseID := vars["course_id"]

	// Логируем полученный ID курса
	log.Printf("Запрашивается курс с ID: %s", courseID)

	// Получаем уроки для курса из базы данных
	lessons, err := getLessonsForCourse(courseID)
	if err != nil {
		http.Error(w, "Ошибка при получении уроков", http.StatusInternalServerError)
		log.Println("Ошибка при получении уроков:", err)
		return
	}

	// Получаем шаблон для рендеринга
	tmpl, ok := TemplateCache["lesson.html"]
	if !ok {
		http.Error(w, "Шаблон не найден", http.StatusInternalServerError)
		return
	}

	// Рендерим шаблон с уроками
	err = tmpl.Execute(w, lessons)
	if err != nil {
		http.Error(w, "Ошибка при отображении страницы", http.StatusInternalServerError)
		log.Println("Ошибка при отображении страницы:", err)
		return
	}
}

// getLessonsForCourse получает все уроки для указанного курса
func getLessonsForCourse(courseID string) ([]Lesson, error) {
	// Запрос к базе данных для получения уроков
	rows, err := database.DB.Query("SELECT id, title, url FROM lessons WHERE course_id = $1", courseID)
	if err != nil {
		log.Println("Ошибка при выполнении запроса:", err)
		return nil, err
	}
	defer rows.Close()

	// Считываем уроки из результата запроса
	var lessons []Lesson
	for rows.Next() {
		var lesson Lesson
		err := rows.Scan(&lesson.ID, &lesson.Title, &lesson.URL)
		if err != nil {
			log.Println("Ошибка при чтении строки результата:", err)
			return nil, err
		}
		lessons = append(lessons, lesson)
	}

	// Возвращаем список уроков
	return lessons, nil
}

func GetLessons() ([]Lesson, error) {
	rows, err := database.DB.Query(`
        SELECT lessons.id, lessons.title, lessons.course_id, courses.title AS course_title, lessons.url
        FROM lessons
        JOIN courses ON lessons.course_id = courses.id
    `)
	if err != nil {
		log.Println("Ошибка при выполнении запроса уроков:", err)
		return nil, err
	}
	defer rows.Close()

	var lessons []Lesson
	for rows.Next() {
		var lesson Lesson
		err := rows.Scan(&lesson.ID, &lesson.Title, &lesson.CourseID, &lesson.CourseTitle, &lesson.URL)
		if err != nil {
			log.Println("Ошибка при чтении строки результата уроков:", err)
			return nil, err
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}
