package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
	title TEXT NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

var DB *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	var install bool
	if os.IsNotExist(err) {
		install = true
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return fmt.Errorf("ошибка создания таблицы: %w", err)
		}
	}

	return nil
}
