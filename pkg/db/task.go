package db

import (
	"database/sql"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int, search string) ([]*Task, error) {
	var (
		rows *sql.Rows
		err  error
	)

	search = strings.TrimSpace(search)

	if search == "" {
		query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`
		rows, err = DB.Query(query, limit)
	} else if t, errTime := time.Parse("02.01.2006", search); errTime == nil {
		searchDate := t.Format("20060102")
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?`
		rows, err = DB.Query(query, searchDate, limit)
	} else {
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE LOWER(title) LIKE LOWER(?) OR LOWER(comment) LIKE LOWER(?) ORDER BY date ASC LIMIT ?`
		likeStr := "%" + search + "%"
		rows, err = DB.Query(query, likeStr, likeStr, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return []*Task{}, err
		}
		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}
