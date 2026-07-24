package db

import (
	"database/sql"
	"os"
	"time"

	"task_scheduler/internal/config"

	_ "modernc.org/sqlite"
)

const schema = `CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(128) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS scheduler_date_idx ON scheduler(date);
`

var db *sql.DB

func Init(cfg *config.DBConfig) error {
	var install bool
	_, err := os.Stat(cfg.Path)
	install = os.IsNotExist(err)

	base, err := sql.Open("sqlite", cfg.Path)
	if err != nil {
		return err
	}

	base.SetMaxIdleConns(1)
	base.SetMaxOpenConns(1)
	base.SetConnMaxIdleTime(time.Minute * 5)
	base.SetConnMaxLifetime(time.Minute * 30)

	if install {
		_, err = base.Exec(schema)
		if err != nil {
			_ = base.Close()
			return err
		}
	}

	db = base
	return nil
}

func Close() {
	if db != nil {
		_ = db.Close()
	}
}
