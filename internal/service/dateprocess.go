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

// сравнение только по дате
func AfterNow(date, now time.Time) bool {
    return date.Format(db.DateFormat) > now.Format(db.DateFormat)
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
        return ruleOfDay(now, date, repeatData)
    
    case "w":
        return ruleOfWeek(now, date, repeatData)

    case "m":
        return ruleOfMonth(now, date, repeatData)

    case "y":
        return ruleOfYear(now, date, repeatData)

    default:
        return "", ErrorUnsupportRepeat
    }
}

func ruleOfDay(now time.Time, date time.Time, repeatData []string) (string, error) {
    if len(repeatData) != 2 {
        return "", fmt.Errorf("%w: invalid days rule", ErrorIncorrectRepeat)
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
}

func ruleOfYear(now time.Time, date time.Time, repeatData []string) (string, error) {
    if len(repeatData) != 1 {
        return "", fmt.Errorf("%w: invalid year rule", ErrorIncorrectRepeat)
    }

    for {
        date = date.AddDate(1, 0, 0)
        if AfterNow(date, now) {
            break
        }
    }
    return date.Format(db.DateFormat), nil
}

func ruleOfWeek(now time.Time, date time.Time, repeatData []string) (string, error) {
    if len(repeatData) != 2 {
        return "", fmt.Errorf("%w: invalid week rule", ErrorIncorrectRepeat)
    }

    daysStr := strings.Split(repeatData[1], ",")
    days := make(map[time.Weekday]struct{})

    for _, d := range daysStr {
        day, err := strconv.Atoi(strings.TrimSpace(d))
        if err != nil || day < 1 || day > 7 {
            return "", fmt.Errorf("%w: invalid weekday value", ErrorIncorrectRepeat)
        }
        
        // Конвертируем воскресенье 7 -> 0, остальные значения совпадают
        weekday := time.Weekday(day % 7)
        days[weekday] = struct{}{}
    }

    current := date.AddDate(0, 0, 1)
    for {
        _, ok := days[current.Weekday()]
        if ok && AfterNow(current, now) {
            return current.Format(db.DateFormat), nil
        }
        current = current.AddDate(0, 0, 1)
    }
}

func ruleOfMonth(now time.Time, date time.Time, repeatData []string) (string, error) {
    if len(repeatData) < 2 {
        return "", fmt.Errorf("%w: invalid month rule", ErrorIncorrectRepeat)
    }

    validDays, err := parseMonthDays(repeatData[1])
    if err != nil {
        return "", err
    }

    var months map[time.Month]struct{}

    if len(repeatData) == 3 {
        months, err = parseMonths(repeatData[2])
        if err != nil {
            return "", err
        }
    }

    limit := date.AddDate(5, 0, 0)
    for current := date.AddDate(0, 0, 1); !current.After(limit); current = current.AddDate(0,0,1){
        if len(months) > 0 {
            if _, ok := months[current.Month()]; !ok {
                continue
            }
        }

        if !AfterNow(current, now) {
            continue
        }

        if isValidMonthDay(current.Day(), daysInMonth(current), validDays) {
            return current.Format(db.DateFormat), nil
        }
    }
    return "", fmt.Errorf("%w: date not found", ErrorIncorrectRepeat)
}

// daysInMonth возвращает количество дней в месяце
func daysInMonth(t time.Time) int {
    return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}

func isValidMonthDay(day, lastDay int, validDays map[int]struct{}) bool {
    if _, ok := validDays[day]; ok {
        return true
    }

    if _, ok := validDays[-1]; ok && day == lastDay {
        return true
    }

    if _, ok := validDays[-2]; ok && day == lastDay-1 {
        return true
    }

    return false
}

func parseMonthDays(s string) (map[int]struct{}, error) {
    days := make(map[int]struct{})
    monthDays := strings.Split(s, ",")

    for _, v := range monthDays {
        day, err := strconv.Atoi(strings.TrimSpace(v))
        if err != nil {
            return nil, fmt.Errorf("%w: invalid month day", ErrorIncorrectRepeat)
        }

        switch {
        case day == -1, day == -2:
            days[day] = struct{}{}

        case day >= 1 && day <= 31:
            days[day] = struct{}{}

        default:
            return nil, fmt.Errorf("%w: invalid month day", ErrorIncorrectRepeat)
        }
    }

    return days, nil
}

func parseMonths(s string) (map[time.Month]struct{}, error) {
    months := make(map[time.Month]struct{})
    monthsStr := strings.Split(s, ",")

    for _, v := range monthsStr {
        month, err := strconv.Atoi(strings.TrimSpace(v))
        if err != nil || month < 1 || month > 12 {
            return nil, fmt.Errorf("%w: invalid month", ErrorIncorrectRepeat)
        }

        m := time.Month(month)
        months[m] = struct{}{}
    }

    return months, nil
}
