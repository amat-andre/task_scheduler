package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"task_scheduler/internal/db"
	"time"
)

var(
	ErrorEmptyRepeat = errors.New("repeat is empty")
	ErrorIncorrectRepeat = errors.New("incorrect repeat format")
	ErrorUnsupportRepeat = errors.New("unsupported repeat format")
	ErrorIncorrectDate = errors.New("incorrect dstart format")
)

const DateFormat = "20060102"

func AfterNow(date, now time.Time) bool {
	return onlyDate(date).After(onlyDate(now)) // true если вызывающий после передаваемого в аргументе (если Data позже)
}

func onlyDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if strings.TrimSpace(repeat) == "" {
		return "", fmt.Errorf("%w: %w", ErrorIncorrectRepeat, ErrorEmptyRepeat)
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrorIncorrectDate, err)
	}

	repeatData := strings.Fields(repeat)
	
	switch repeatData[0]{
	case "d":
		if len(repeatData) != 2 {
			return "", fmt.Errorf("%w: invalid days value", ErrorIncorrectRepeat)
		}

		days, err := strconv.Atoi(repeatData[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("%w: invalid days interval", ErrorIncorrectRepeat)
		}

		for {
    		date = date.AddDate(0, 0, days)
   			if AfterNow(date, now) {
        		break
    		}
		} 
		return date.Format(DateFormat), nil

	case "y":
		if len(repeatData) != 1 {
			return "", fmt.Errorf("%w: invalid year value", ErrorIncorrectRepeat)
		}
		
		for {
    		date = date.AddDate(1, 0, 0)
   			if AfterNow(date, now) {
        		break
    		}
		} 
		return date.Format(DateFormat), nil

	default:
		return "", ErrorUnsupportRepeat
	}
}

func CheckDate(task *db.Task) error {
	now := time.Now()

	if strings.TrimSpace(task.Date) == "" {
		task.Date = now.Format(DateFormat)
	} 

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return ErrorIncorrectDate
	}

	var nextDate string
	if strings.TrimSpace(task.Repeat) != "" {
		nextDate, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// если сегодня (now) больше task.Date (t)
	if AfterNow(now, t) {
        if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
            task.Date = now.Format(DateFormat)
        } else {
            // в противном случае, берём вычисленную ранее следующую дату
            task.Date = nextDate
        }
    } 
	return nil
}