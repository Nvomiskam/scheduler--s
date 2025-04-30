// Пакет db реализует БД, используемую для хранения задач на SQLite.
// Путь к БД берется из TODO_DBFILE или используется scheduler.db по умолчанию
package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

type DB struct {
	db *sql.DB
}

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
func NewDB(dbFile string) (*DB, error) {

	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		if _, err := conn.Exec(schema); err != nil {
			return nil, err
		}
	}
	return &DB{db: conn}, nil
}

// GetDBPath возвращает путь к файлу БД из переменной окружения или по умолчанию
func GetDBPath() string {
	if path := os.Getenv("TODO_DBFILE"); path != "" {
		return path
	}
	return "scheduler.db"
}
