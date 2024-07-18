package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB инициализирует базу данных
func InitDB() {
	// Получение пути к файлу базы данных из переменной окружения
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		// Использование текущей рабочей директории
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dbFile = filepath.Join(cwd, "scheduler.db")
	}

	log.Printf("Using database file: %s", dbFile)

	var install bool
	if _, err := os.Stat(dbFile); err != nil {
		install = true
	}

	// Открытие базы данных
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	DB = db

	if install {
		log.Println("Creating new database...")

		// Создание таблицы и индекса, если база данных новая
		createTableSQL := `CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date TEXT NOT NULL CHECK(length(date) = 8),
            title TEXT NOT NULL,
            comment TEXT,
            repeat TEXT CHECK(length(repeat) <= 128)
        );`
		_, err := DB.Exec(createTableSQL)
		if err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}

		// Создание индекса по полю date
		createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`
		_, err = DB.Exec(createIndexSQL)
		if err != nil {
			log.Fatalf("Failed to create index: %v", err)
		}

		log.Println("Database created successfully.")
	} else {
		log.Println("Database already exists.")
	}
}

// AddTask добавляет задачу в базу данных и возвращает идентификатор новой задачи
func AddTask(date, title, comment, repeat string) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, date, title, comment, repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
