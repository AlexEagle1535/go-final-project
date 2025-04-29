package api

import (
	"net/http"

	"github.com/AlexEagle1535/go-final-project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJSON(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}
