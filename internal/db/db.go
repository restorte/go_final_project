package db

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var DB *sqlx.DB

type Task struct {
	ID      int64  `db:"id" json:"id,string,omitempty"`
	Date    string `db:"date" json:"date"`
	Title   string `db:"title" json:"title"`
	Comment string `db:"comment" json:"comment,omitempty"`
	Repeat  string `db:"repeat" json:"repeat,omitempty"`
}

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	needCreate := err != nil

	db, err := sqlx.Connect("sqlite", dbFile)
	if err != nil {
		return err
	}
	DB = db

	if needCreate {
		schema := `
		CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT '',
			title VARCHAR(255) NOT NULL,
			comment TEXT,
			repeat VARCHAR(128)
		);
		CREATE INDEX idx_date ON scheduler(date);
		`
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int, search string) ([]Task, error) {
	var query string
	var args []interface{}

	if search == "" {
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		args = []interface{}{limit}
	} else if isDate(search) {
		date, err := parseDate(search)
		if err != nil {
			return tasksByText(limit, search)
		}
		query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
		args = []interface{}{date, limit}
	} else {
		return tasksByText(limit, search)
	}

	var tasks []Task
	err := DB.Select(&tasks, query, args...)
	if err != nil {
		return nil, err
	}
	if tasks == nil {
		return []Task{}, nil
	}
	return tasks, nil
}

func isDate(s string) bool {
	if len(s) != 10 {
		return false
	}
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return false
	}
	for _, p := range parts {
		if _, err := strconv.Atoi(p); err != nil {
			return false
		}
	}
	return true
}

func parseDate(s string) (string, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return "", errors.New("invalid date format")
	}
	day := parts[0]
	month := parts[1]
	year := parts[2]
	if len(day) != 2 || len(month) != 2 || len(year) != 4 {
		return "", errors.New("invalid date format")
	}
	return year + month + day, nil
}

func tasksByText(limit int, search string) ([]Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler 
	          WHERE LOWER(title) LIKE LOWER(?) OR LOWER(comment) LIKE LOWER(?) 
	          ORDER BY date LIMIT ?`
	like := "%" + search + "%"
	var tasks []Task
	err := DB.Select(&tasks, query, like, like, limit)
	if err != nil {
		return nil, err
	}
	if tasks == nil {
		return []Task{}, nil
	}
	return tasks, nil
}
