package main //m@8246.ru

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/1ssk/MicroHack.git/internal/database"
	"github.com/1ssk/MicroHack.git/internal/handlers"
	"github.com/1ssk/MicroHack.git/internal/middleware"
	"github.com/gorilla/mux"
)

func main() {

	middleware.SetupLogging()

	if err := database.Init(); err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	defer database.Close()

	// Список шаблонов для кэширования
	templates := []string{"register.html", "admin.html", "login.html", "index.html", "lesson.html", "video.html", "change_password.html"}

	// Инициализация кэша шаблонов
	handlers.InitTemplateCache("web/templates", templates)

	// Чтение порта из переменных окружения
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Порт по умолчанию
	}

	// Настройка роутера
	router := mux.NewRouter()

	// Подключение middleware
	router.Use(middleware.LoggingMiddleware)

	// Маршруты не защищенные
	router.HandleFunc("/login", handlers.LoginHandler).Methods("GET", "POST")
	// router.HandleFunc("/admin/create-course", handlers.LoginHandler).Methods("GET", "POST")
	// router.HandleFunc("/admin/assign-course", handlers.LoginHandler).Methods("GET", "POST")
	// router.HandleFunc("/admin/create-lesson", handlers.LoginHandler).Methods("GET", "POST")
	// router.HandleFunc("/api/user/{user_id}/courses", handlers.LoginHandler).Methods("GET", "POST")
	// router.HandleFunc("/register", handlers.RegisterHandler).Methods("GET", "POST")
	// router.HandleFunc("/api/my-courses", handlers.MyCoursesHandler).Methods("GET")
	// router.HandleFunc("/lessons/{course_id}", handlers.LessonHandler).Methods("GET")
	// router.HandleFunc("/video/{lesson_id}", handlers.VideoHandler).Methods("GET")

	// Защищенные маршруты с авторизацией
	router.Handle("/admin/change-password", middleware.AuthMiddleware(http.HandlerFunc(handlers.ChangePasswordHandler))).Methods("GET", "POST")
	router.Handle("/", middleware.AuthMiddleware(http.HandlerFunc(handlers.IndexHandler))).Methods("GET", "POST")
	router.Handle("/admin", middleware.AuthMiddleware(http.HandlerFunc(handlers.AdminHandler))).Methods("GET", "POST")
	router.Handle("/admin/create-course", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateCourse))).Methods("POST")
	router.Handle("/admin/assign-course", middleware.AuthMiddleware(http.HandlerFunc(handlers.AssignCourse))).Methods("POST")
	router.Handle("/admin/create-lesson", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateLesson))).Methods("POST")
	router.Handle("/api/user/{user_id}/courses", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserCourses))).Methods("GET")
	router.Handle("/register", middleware.AuthMiddleware(http.HandlerFunc(handlers.RegisterHandler))).Methods("GET", "POST")
	router.Handle("/api/my-courses", middleware.AuthMiddleware(http.HandlerFunc(handlers.MyCoursesHandler))).Methods("GET")
	router.Handle("/lessons/{course_id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.LessonHandler))).Methods("GET")
	router.Handle("/video/{lesson_id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.VideoHandler))).Methods("GET")
	router.Handle("/admin/delete-course", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteCourse))).Methods("POST")
	router.Handle("/admin/delete-user", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteUser))).Methods("POST")
	router.Handle("/admin/delete-lesson", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteLesson))).Methods("POST")

	// Маршрут для выхода (logout)
	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		// Удаляем куки с токеном
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   "",
			Expires: time.Now().Add(-time.Hour), // Устанавливаем время в прошлом для удаления куки
		})
		http.Redirect(w, r, "/login", http.StatusFound)
	}).Methods("GET")

	// Обработчик статических файлов
	router.PathPrefix("/web/static/").Handler(http.StripPrefix("/web/static/", http.FileServer(http.Dir("./web/static"))))

	// Создание HTTP-сервера
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Канал для обработки сигналов остановки
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill) // SIGINT, SIGTERM

	// Запуск сервера в отдельной горутине
	go func() {
		log.Printf("Сервер запущен на порту :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Блокировка до получения сигнала остановки
	<-stop
	log.Println("Получен сигнал остановки сервера")

	// Контекст с тайм-аутом для завершения работы
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка завершения сервера: %v", err)
	}
	log.Println("Сервер корректно завершил работу")
}
