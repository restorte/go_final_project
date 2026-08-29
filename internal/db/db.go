package db

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var DB *sqlx.DB

type Task struct {
	ID      int64  `db:"id"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
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
