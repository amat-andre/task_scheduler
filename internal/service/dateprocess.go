package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"task_scheduler/internal/db"
	"time"
)

var (
	ErrorEmptyRepeat     = errors.New("repeat is empty")
	ErrorIncorrectRepeat = errors.New("incorrect repeat format")
	ErrorUnsupportRepeat = errors.New("unsupported repeat format")
	ErrorIncorrectStart  = errors.New("incorrect dstart format")
	ErrorIncorrectDate   = errors.New("incorrect date format")
)

func AfterNow(date, now time.Time) bool {
	return onlyDate(date).After(onlyDate(now))
}

func onlyDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if strings.TrimSpace(repeat) == "" {
		return "", fmt.Errorf("%w: %w", ErrorIncorrectRepeat, ErrorEmptyRepeat)
	}

	date, err := time.Parse(db.DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrorIncorrectStart, err)
	}

	repeatData := strings.Fields(repeat)

	switch repeatData[0] {
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
		return date.Format(db.DateFormat), nil

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
		return date.Format(db.DateFormat), nil

	default:
		return "", ErrorUnsupportRepeat
	}
}

func CheckDate(task *db.Task) error {
	now := time.Now()

	if strings.TrimSpace(task.Date) == "" {
		task.Date = now.Format(db.DateFormat)
	}

	t, err := time.Parse(db.DateFormat, task.Date)
	if err != nil {
		return ErrorIncorrectDate
	}

	var nextDate string
	if strings.TrimSpace(task.Repeat) != "" {
		nextDate, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("check date failed: %w", err)
		}
	}

	if AfterNow(now, t) {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format(db.DateFormat)
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			task.Date = nextDate
		}
	}
	return nil
}
