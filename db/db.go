package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// Инициализируем базу данных
func InitDB() {
	// Получение пути к файлу базы данных из переменной окружения
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		// Использование текущей рабочей директории
		dbFile = "scheduler.db"
	}

	log.Printf("Using database file: %s", dbFile)

	var install bool
	if _, err := os.Stat(dbFile); err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			log.Fatalf("Failed to check if database file exists: %v", err)
		}
	}

	// Открытие базы данных
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	DB = db

	if !install {
		log.Println("Database already exists.")
		return
	}

	log.Println("Creating new database...")

	// Создание таблицы и индекса, если база данных новая
	createTableSQL := `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL CHECK(length(date) = 8),
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT CHECK(length(repeat) <= 128)
	);`
	_, err = DB.Exec(createTableSQL)
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
}

// Добавляем задачу в базу данных и возвращаем идентификатор новой задачи
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

// Структура задачи
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Возвращаем список ближайших задач из базы данных
// В задании этого нет, но если фронтенд будет поддерживать пагинацию, то это пригодится
func GetTasks(limit, offset int) ([]Task, error) {
	now := time.Now().Format("20060102")
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date >= ? ORDER BY date LIMIT ? OFFSET ?`
	rows, err := DB.Query(query, now, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var id int64
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		task.ID = fmt.Sprintf("%d", id)
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// Возвращаем задачи по заданной дате
func GetTasksByDate(date string, limit, offset int) ([]Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ? OFFSET ?`
	rows, err := DB.Query(query, date, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var id int64
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		task.ID = fmt.Sprintf("%d", id)
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// Выполняем поиск задач по подстроке в заголовке или комментарии
func SearchTasks(search string, limit, offset int) ([]Task, error) {
	searchTerm := "%" + search + "%"
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ? OFFSET ?`
	rows, err := DB.Query(query, searchTerm, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var id int64
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		task.ID = fmt.Sprintf("%d", id)
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// Возвращаем задачу по её идентификатору
func GetTaskByID(id int64) (Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := DB.QueryRow(query, id)

	var task Task
	var taskID int64
	err := row.Scan(&taskID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, fmt.Errorf("задача не найдена")
		}
		return Task{}, err
	}
	task.ID = fmt.Sprintf("%d", taskID)
	return task, nil
}

// Обновляем задачу в базе данных
func UpdateTask(id, date, title, comment, repeat string) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, date, title, comment, repeat, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// Удаляем задачу из базы данных
func DeleteTask(id int64) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}
