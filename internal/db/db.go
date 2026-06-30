package db

import(
	"time"
	"os"

	"database/sql"
	_ "modernc.org/sqlite"
)

const (
schema = `CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(128) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX scheduler_date_idx ON scheduler(date);
`

defFileDB = "scheduler.db"
)

var db *sql.DB

func Init() error {	
	dbFile := getFileDB()
	var install bool
	_, err := os.Stat(dbFile)
	install = os.IsNotExist(err)


	base, err := sql.Open("sqlite", dbFile)
    if err != nil {
        return err
    }

	base.SetMaxIdleConns(2)
    base.SetMaxOpenConns(4)
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

func Close(){
	if db != nil {
		_ = db.Close()
	}
}

func getFileDB() string {
	if dbFile := os.Getenv("TODO_DBFILE"); len(dbFile) > 0 {
		return dbFile
	}
	return defFileDB
}