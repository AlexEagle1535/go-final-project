package api

import (
	"net/http"

	"github.com/AlexEagle1535/go-final-project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID is required"})
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete task"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}
