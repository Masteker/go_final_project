package handlers

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/Masteker/go_final_project/database"
	"net/http"
)

func HandleDeleteTask(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
			return
		}

		err := database.DeleteTask(db, id)
		if err != nil {
			if err.Error() == "задача не найдена" {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			}
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}