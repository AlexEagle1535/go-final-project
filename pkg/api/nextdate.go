package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const Tformat = "20060102"

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if now == "" || date == "" || repeat == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	nowTime, err := time.Parse(Tformat, now)
	if err != nil {
		http.Error(w, "invalid now date format", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}
	date, err := time.Parse(Tformat, dstart)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга даты: %w", err)
	}

	repeatSplit := strings.Split(repeat, " ")

	switch repeatSplit[0] {
	case "d":
		if len(repeatSplit) != 2 {
			return "", fmt.Errorf("invalid 'd' format")
		}
		interval, err := strconv.Atoi(repeatSplit[1])
		if err != nil || interval <= 0 || interval > 400 {
			return "", fmt.Errorf("invalid day interval")
		}
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(Tformat), nil

	case "y":
		if len(repeatSplit) != 1 {
			return "", fmt.Errorf("invalid 'y' format")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(Tformat), nil

	case "w":
		if len(repeatSplit) != 2 {
			return "", fmt.Errorf("invalid 'w' format")
		}
		days := [7]bool{}
		dayList := strings.Split(repeatSplit[1], ",")
		for _, d := range dayList {
			dayNum, err := strconv.Atoi(d)
			if err != nil || dayNum < 1 || dayNum > 7 {
				return "", fmt.Errorf("invalid weekday: %s", d)
			}
			days[dayNum-1] = true
		}
		for {
			date = date.AddDate(0, 0, 1)
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			if days[weekday-1] && afterNow(date, now) {
				break
			}
		}
		return date.Format(Tformat), nil

	case "m":
		if len(repeatSplit) < 2 {
			return "", fmt.Errorf("invalid 'm' format")
		}
		dayFlags := [32]bool{}
		monthFlags := [13]bool{}
		var last, prelast, backLast, backPrelast bool
		// Заполняем допустимые дни
		daysPart := strings.Split(repeatSplit[1], ",")
		for _, d := range daysPart {
			dayNum, err := strconv.Atoi(d)
			if err != nil || (dayNum < 1 && dayNum != -1 && dayNum != -2) || dayNum > 31 {
				return "", fmt.Errorf("invalid day in 'm' rule: %s", d)
			}
			switch dayNum {
			case -1:
				last = true
			case -2:
				prelast = true
			default:
				dayFlags[dayNum] = true
			}
		}

		// Заполняем допустимые месяцы
		if len(repeatSplit) >= 3 {
			monthsPart := strings.Split(repeatSplit[2], ",")
			for _, m := range monthsPart {
				monthNum, err := strconv.Atoi(m)
				if err != nil || monthNum < 1 || monthNum > 12 {
					return "", fmt.Errorf("invalid month in 'm' rule: %s", m)
				}
				monthFlags[monthNum] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				monthFlags[i] = true
			}
		}

		for {
			year, month, _ := date.Date()
			lastDay := lastDayOfMonth(date)
			monthInt := int(month)
			if !monthFlags[monthInt] {
				// Месяц не разрешён — переходим к следующему месяцу
				date = time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
				continue
			}
			if last && !dayFlags[lastDay] {
				dayFlags[lastDay] = true
				backLast = true
			}
			if prelast && !dayFlags[lastDay-1] {
				dayFlags[lastDay-1] = true
				backPrelast = true
			}
			// Ищем день в месяце
			for day := 1; day <= lastDay; day++ {
				if !dayFlags[day] {
					continue
				}
				if day > lastDay {
					// В месяце нет такого дня (например, 31-го)
					continue
				}
				candidate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
				if afterNow(candidate, now) {
					return candidate.Format(Tformat), nil
				}
			}
			if backLast {
				dayFlags[lastDay] = false
				backLast = false
			}
			if backPrelast {
				dayFlags[lastDay-1] = false
				backPrelast = false
			}
			// Ничего не подошло — переход на следующий месяц
			date = time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
		}
	default:
		return "", fmt.Errorf("unsupported repeat format: %s", repeatSplit[0])
	}
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
