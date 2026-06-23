package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"
)

var ErrorTaskNotFound = errors.New("identifier task not found") // возможно слово идентификатор нужно будет удалить, посмотреть логику

type Task struct {
	ID 		string	`json:"id"`
    Date    string  `json:"date"`
    Title   string 	`json:"title"`
    Comment string  `json:"comment"`
    Repeat  string 	`json:"repeat"`
}

func Tasks(limit int) ([]*Task, error){
	query := "SELECT * FROM scheduler ORDER BY date LIMIT ?"

	rows, err := db.Query(query, limit)
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
	/*
	еще можно сделать так, надо узнать как лучше

	for rows.Next() {
		var task Task
		var id int64
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, task)
	}
	*/
	if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("presence errors in rows: %w", err)
    }

	return tasks, nil
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
        return 0, err
    }

	id, err := res.LastInsertId()
    return id, err
}

func GetTask(id string) (*Task, error){ // здесь возможно тоже придется заменить id на число
	task := &Task{}

	query := "SELECT * FROM scheduler WHERE id = :id"
	row := db.QueryRow(query, sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorTaskNotFound
		}
		return nil, err
	}
	
	return task, nil
}

func UpdateTask(task *Task) error {
    query := "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id"
    res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
    if err != nil {
        return fmt.Errorf("failed to update: %w", err)
    }
    
    count, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to update: %w", err)
    }
    if count == 0 {
        return ErrorTaskNotFound
    }
    return nil
}

func UpdateDate(nextDate string, id string) error {
	query := "UPDATE scheduler SET date = :date WHERE id = :id"
    res, err := db.Exec(query,
		sql.Named("date", nextDate),
		sql.Named("id", id))
    if err != nil {
        return fmt.Errorf("failed to update date: %w", err)
    }
    
    count, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to update date: %w", err)
    }
    if count == 0 {
        return ErrorTaskNotFound
    }
    return nil
}

func DeleteTask(id string) error{
	query := "DELETE FROM scheduler WHERE id = :id"
	res, err := db.Exec(query, sql.Named("id", id))
    if err != nil {
        return fmt.Errorf("failed to delete: %w", err)
	}

	count, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to delete: %w", err)
    }
    if count == 0 {
        return ErrorTaskNotFound
    }
    return nil
}