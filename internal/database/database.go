package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Подключение PostgreSQL драйвера
)

var DB *sql.DB

// Init инициализирует подключение к базе данных
func Init() error {
	// Загружаем переменные из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке файла .env")
	}

	// Чтение переменных окружения для подключения к базе данных
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Проверяем, что все необходимые переменные заданы
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		return fmt.Errorf("некоторые переменные окружения не заданы")
	}

	// Формируем строку подключения
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", dbUser, dbPassword, dbHost, dbPort, dbName)
	log.Printf("Попытка подключения: %s", connStr)

	// Открываем соединение с базой данных
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("Ошибка при подключении к базе данных: %v", err)
	}

	// Проверяем соединение с базой
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("Ошибка при проверке подключения: %v", err)
	}

	log.Println("Подключение к базе данных успешно!")
	return nil
}

// Close закрывает соединение с базой данных
func Close() {
	if err := DB.Close(); err != nil {
		log.Printf("Ошибка при закрытии соединения с базой данных: %v", err)
	}
}
