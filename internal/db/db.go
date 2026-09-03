package db

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var DB *sqlx.DB

type Task struct {
	ID      int64  `db:"id" json:"id,omitempty"`
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
