// Пакет db реализует БД, используемую для хранения задач на SQLite.
// Путь к БД берется из TODO_DBFILE или используется scheduler.db по умолчанию
package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

// Используется для глобального подключения к базе данных
var db *sql.DB

// Константа с SQL командами для создания таблицы scheduler и индекса по колонке date
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT,
    repeat VARCHAR(128) DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`

// Init инициализирует подключение к БД
// dbFile - путь к файлу базы данных
func Init(dbFile string) error {

	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		if _, err := db.Exec(schema); err != nil {
			return err
		}
	}
	return nil
}

// GetDB возвращает глобальное подключение к БД
func GetDB() *sql.DB {
	return db
}

// GetDBPath возвращает путь к файлу БД из переменной окружения или по умолчанию
func GetDBPath() string {
	if path := os.Getenv("TODO_DBFILE"); path != "" {
		return path
	}
	return "scheduler.db"
}
