// Пакет db реализует БД, используемую для хранения задач на SQLite.
package db

import (
	"database/sql"
	"fmt"
	"time"

	"go1f/pkg/rules"
)

// Структура Task представляет задачу в планировщике
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в базу данных
func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления задачи в базу: %v", err)
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения id последней добавленной задачи: %v", err)
	}
	return id, nil
}

// Tasks возвращает список задач, отсортированных по дате
func Tasks(limit int) ([]*Task, error) {
	rows, err := db.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе данных: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения данных: %w", err)
		}
		tasks = append(tasks, &t)
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// GetFilteredTasks позволяет пользоваться поиском задач по указанному слову или по дате
func GetFilteredTasks(search string, limit int) ([]*Task, error) {
	var rows *sql.Rows
	var err error

	if date, errParse := time.Parse("02.01.2006", search); errParse == nil {
		dateStr := date.Format(rules.DateFormat)
		query := `SELECT * FROM scheduler WHERE date = ? LIMIT ?`
		rows, err = db.Query(query,
			sql.Named("date", dateStr),
			sql.Named("limit", limit))
	} else {
		like := "%" + search + "%"
		query := `SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit`
		rows, err = db.Query(query,
			sql.Named("search", like),
			sql.Named("limit", limit))
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе данных: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("ошибка чтения данных: %w", err)
		}
		tasks = append(tasks, &t)
	}

	return tasks, nil
}

// GetTask возвращает задачу по её идентификатору
func GetTask(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := db.QueryRow(query, id)

	var t Task
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения данных: %w", err)
	}

	return &t, nil
}

// UpdateTask обновляет параметры задачи по идентификатору
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновлённых строк: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки количества удаленных строк: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("задача с id=%s не найдена", id)
	}

	return nil
}

// UpdateDate обновляет дату задачи по её идентификатору
func UpdateDate(nextDate string, id string) error {
	_, err := db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
	if err != nil {
		return fmt.Errorf("failed to update task date: %v", err)
	}
	return nil
}
