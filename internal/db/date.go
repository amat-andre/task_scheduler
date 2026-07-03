package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func UpdateDate(nextDate string, id string) error {
	query := `UPDATE scheduler 
        SET date = :date 
        WHERE id = :id`
        
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