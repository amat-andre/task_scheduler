package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

var ErrorTaskNotFound = errors.New("identifier task not found")

type Task struct {
	ID 		string	`json:"id"`
    Date    string  `json:"date"`
    Title   string 	`json:"title"`
    Comment string  `json:"comment"`
    Repeat  string 	`json:"repeat"`
}

func AddTask(task *Task) (int64, error){
	query := `INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES (:date, :title, :comment, :repeat)`

    res, err := db.Exec(query, 
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
    if err != nil {	
        return 0, fmt.Errorf("failed to add task: %w", err)
    }

	id, err := res.LastInsertId()
    return id, err
}

func GetTask(id string) (*Task, error){
	task := &Task{}

	query := `SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE id = :id`

	row := db.QueryRow(query, sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	
	return task, nil
}

func UpdateTask(task *Task) error {
    query := `UPDATE scheduler 
		SET date = :date, title = :title, comment = :comment, repeat = :repeat 
		WHERE id = :id`

    res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
    if err != nil {
        return fmt.Errorf("failed to update task: %w", err)
    }
    
    count, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to update task: %w", err)
    }
    if count == 0 {
        return ErrorTaskNotFound
    }
    return nil
}

func DeleteTask(id string) error{
	query := `DELETE FROM scheduler WHERE id = :id`

	res, err := db.Exec(query, sql.Named("id", id))
    if err != nil {
        return fmt.Errorf("failed to delete task: %w", err)
	}

	count, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to delete task: %w", err)
    }
    if count == 0 {
        return ErrorTaskNotFound
    }
    return nil
}


func Tasks(search string, limit int) ([]*Task, error){
	var (
		query string
		err error
		rows *sql.Rows
	)

	if search == "" {
		query = `SELECT id, date, title, comment, repeat 
			FROM scheduler 
			ORDER BY date 
			LIMIT :limit`

		rows, err = db.Query(query, sql.Named("limit", limit))
	} else {
		date, err := time.Parse("02.01.2006", search)

		if err == nil {
			dateQuery := date.Format("20060102")
			
			query = `SELECT id, date, title, comment, repeat 
				FROM scheduler 
				WHERE date = :date 
				LIMIT :limit`

			rows, err = db.Query(query, 
				sql.Named("date", dateQuery), 
				sql.Named("limit", limit))
		} else {
			query = `SELECT id, date, title, comment, repeat 
				FROM scheduler 
				WHERE title LIKE :search OR comment LIKE :search 
				ORDER BY date LIMIT :limit`

			rows, err = db.Query(query, 
				sql.Named("search", "%"+search+"%"), 
				sql.Named("limit", limit))
		}
	}

	if err != nil {
    	return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	tasks := []*Task{}
	for rows.Next(){
		task := &Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
    		return nil, fmt.Errorf("failed to scan: %w", err)
		}
		tasks = append(tasks, task)
	}
	
	if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("presence errors in rows: %w", err)
    }

	return tasks, nil
}