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
func (d *DB) AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := d.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("failed to add task to the database: %v", err)
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last inserted task ID: %v", err)
	}
	return id, nil
}

// Tasks возвращает список задач, отсортированных по дате
func (d *DB) Tasks(limit int) ([]*Task, error) {
	rows, err := d.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query the database: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, fmt.Errorf("failed to read task data: %w", err)
		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error processing rows: %w", err)
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// GetFilteredTasks позволяет пользоваться поиском задач по указанному слову или по дате
func (d *DB) GetFilteredTasks(search string, limit int) ([]*Task, error) {
	var rows *sql.Rows
	var err error

	if date, errParse := time.Parse("02.01.2006", search); errParse == nil {
		dateStr := date.Format(rules.DateFormat)
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?`
		rows, err = d.db.Query(query,
			sql.Named("date", dateStr),
			sql.Named("limit", limit))
	} else {
		like := "%" + search + "%"
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit`
		rows, err = d.db.Query(query,
			sql.Named("search", like),
			sql.Named("limit", limit))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query the database: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("failed to read task data: %w", err)
		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error processing rows: %w", err)
	}

	return tasks, nil
}

// GetTask возвращает задачу по её идентификатору
func (d *DB) GetTask(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := d.db.QueryRow(query, id)

	var t Task
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, fmt.Errorf("failed to read task data: %w", err)
	}

	return &t, nil
}

// UpdateTask обновляет параметры задачи по идентификатору
func (d *DB) UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := d.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve the number of updated rows: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (d *DB) DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check the number of deleted rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with id=%s not found", id)
	}

	return nil
}

// UpdateDate обновляет дату задачи по её идентификатору
func (d *DB) UpdateDate(nextDate string, id string) error {
	_, err := d.db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
	if err != nil {
		return fmt.Errorf("failed to update task date: %v", err)
	}
	return nil
}
