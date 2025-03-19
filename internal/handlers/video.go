package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/1ssk/MicroHack.git/internal/database"

	"github.com/gorilla/mux"
)

// VideoHandler обрабатывает запросы для отображения видеоурока
func VideoHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID урока из URL
	vars := mux.Vars(r)
	lessonID := vars["lesson_id"]

	// Получаем информацию о видеоуроке
	lesson, err := getLessonByID(lessonID)
	if err != nil {
		http.Error(w, "Ошибка при получении урока", http.StatusInternalServerError)
		log.Println("Ошибка при получении урока:", err)
		return
	}

	// Проверяем, существует ли кэшированный шаблон
	cacheMutex.RLock()
	tmpl, ok := TemplateCache["video.html"]
	cacheMutex.RUnlock()
	if !ok {
		// Если шаблон не найден в кэше, загружаем его
		tmplPath := "templates/video.html"
		tmpl, err = template.ParseFiles(tmplPath)
		if err != nil {
			http.Error(w, "Ошибка при обработке шаблона", http.StatusInternalServerError)
			log.Println("Ошибка при обработке шаблона:", err)
			return
		}
		// Добавляем шаблон в кэш
		cacheMutex.Lock()
		TemplateCache["video.html"] = tmpl
		cacheMutex.Unlock()
	}

	// Рендерим шаблон с данными урока
	err = tmpl.Execute(w, lesson)
	if err != nil {
		http.Error(w, "Ошибка при отображении страницы", http.StatusInternalServerError)
		log.Println("Ошибка при отображении страницы:", err)
		return
	}
}

// getLessonByID получает информацию о конкретном уроке по ID
func getLessonByID(lessonID string) (*Lesson, error) {
	// Запрос к базе данных для получения информации о видеоуроке
	row := database.DB.QueryRow("SELECT id, title, url FROM lessons WHERE id = $1", lessonID)

	// Инициализируем переменную для хранения данных урока
	var lesson Lesson
	err := row.Scan(&lesson.ID, &lesson.Title, &lesson.URL)
	if err != nil {
		log.Println("Ошибка при получении данных урока:", err)
		return nil, err
	}

	// Возвращаем информацию о уроке
	return &lesson, nil
}
