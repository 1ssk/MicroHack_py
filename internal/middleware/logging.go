package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Функция для настройки логирования (вывод в файл и консоль)
func SetupLogging() {
	// Создаём директорию для логов, если её нет
	if err := os.MkdirAll("log", os.ModePerm); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Имя файла с логами на основе текущей даты
	logFileName := fmt.Sprintf("log/server-%s.log", time.Now().Format("2006-01-02"))
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Настроим мульти-вывод: в консоль и в файл
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.Println("Logging initialized")
}

type logEntry struct {
	Method     string        `json:"method"`
	Path       string        `json:"path"`
	RemoteAddr string        `json:"remote_addr"`
	StatusCode int           `json:"status_code"`
	Duration   time.Duration `json:"duration"`
}

// Middleware для логирования запросов
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Временная запись, чтобы перехватить статус ответа
		recorder := &responseRecorder{ResponseWriter: w}
		next.ServeHTTP(recorder, r)

		duration := time.Since(startTime)
		entry := logEntry{
			Method:     r.Method,
			Path:       r.URL.Path,
			RemoteAddr: r.RemoteAddr,
			StatusCode: recorder.statusCode,
			Duration:   duration,
		}
		logJSON, err := json.Marshal(entry)
		if err != nil {
			log.Printf("Ошибка форматирования лога: %v", err)
			return
		}
		log.Println(string(logJSON))
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	// Если статус код еще не установлен, устанавливаем его
	if rr.statusCode == 0 {
		rr.statusCode = code
	}
	rr.ResponseWriter.WriteHeader(code)
}
