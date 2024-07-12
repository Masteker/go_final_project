package handlers

import (
	"encoding/json"
	"fmt"
	"go_final_project/internal/models&service/store/task"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"
)

// Task обрабатывает HTTP-запросы для выполнения CRUD-операций над задачами в приложении-планировщике.
// Он поддерживает методы GET, POST, PUT и DELETE.
//
// GET:
//   - Извлекает задачу по её ID из базы данных.
//   - Возвращает JSON-объект со всеми полями задачи.
//   - Если ID пуст или не найден, возвращает ошибку 400 Bad Request.
//
// POST:
//   - Создает новую задачу в базе данных.
//   - Возвращает JSON-объект с ID созданной задачи.
//   - Если тело запроса недействительно или отсутствуют поля, возвращает ошибку 400 Bad Request.
//
// PUT:
//   - Обновляет существующую задачу в базе данных.
//   - Возвращает пустой JSON-объект.
//   - Если тело запроса недействительно или отсутствуют поля, возвращает ошибку 400 Bad Request.
//   - Если указанный ID не найден в базе данных, возвращает ошибку 400 Bad Request.
//
// DELETE:
//   - Удаляет задачу из базы данных по её ID.
//   - Возвращает пустой JSON-объект.
//   - Если ID пуст или не найден, возвращает ошибку 400 Bad Request.
//
// По умолчанию, если метод запроса не GET, POST, PUT или DELETE, возвращает ошибку 500 Internal Server Error.
func (h *Handler) Task(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	switch r.Method {
	case http.MethodGet:
		id, err := h.GetID(r)
		if err != nil {
			h.SendErr(w, err, http.StatusBadRequest)
			return
		}
		task, err := h.service.Store.GetTask(id)
		if err != nil {
			h.SendErr(w, err, http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(task)
		if err != nil {
			h.SendErr(w, err, http.StatusInternalServerError)
			return
		}
		h.logger.Infof("sent response via handler Task (method %s)", r.Method)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(response)
	case http.MethodPut:
		var task task.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			err = fmt.Errorf("can't parse response")
			h.SendErr(w, err, http.StatusBadRequest)
			return
		}
		err = task.CheckTask()
		if err != nil {
			h.SendErr(w, err, http.StatusBadRequest)
			return
		} else {
			id, err := strconv.Atoi(task.ID)
			if err != nil {
				err = fmt.Errorf("cat not parse id")
				h.SendErr(w, err, http.StatusBadRequest)
				return
			}
			err = h.service.Store.CheckID(id)
			if err != nil {
				h.SendErr(w, err, http.StatusBadRequest)
				return
			}
			task.Date, err = task.CheckDate()
			if err != nil {
				h.SendErr(w, err, http.StatusBadRequest)
				return
			}
		}
		err = h.service.Store.Update(&task)
		if err != nil {
			h.SendErr(w, err, http.StatusInternalServerError)
			return
		}
		h.logger.Infof("sent response via handler Task (method %s)", r.Method)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprint(w, "{}")
	case http.MethodPost:
		var id struct {
			ID int `json:"id"`
		}
		var request task.Task
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			h.SendErr(w, err, http.StatusBadRequest)
			return
		}

		err = request.CheckTask()
		if err != nil {
			h.SendErr(w, err, http.StatusBadRequest)
			return
		}
		nextDate, err := request.CompleteRequest()
		if err != nil {
			h.SendErr(w, err, http.StatusBadRequest)
			return
		}
		request.Date = nextDate
		if request.Date != "" {
			nextDate, err = request.CheckDate()
			if err != nil {
				h.SendErr(w, err, http.StatusBadRequest)
			}
			request.Date = nextDate
		}
		lastInsertID, err := h.service.Store.Insert(&request)
		if err != nil {
			h.SendErr(w, err, http.StatusInternalServerError)
			return
		}
		id.ID = lastInsertID
		response, err := json.Marshal(id)
		if err != nil {
			h.SendErr(w, err, http.StatusInternalServerError)
			return
		}
		h.logger.Infof("sent response via handler Task (method %s)", r.Method)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(response)
	case http.MethodDelete:
		id, err := h.GetID(r)
		if err != nil {
			h.SendErr(w, err, http.StatusBadRequest)
			return
		}
		// TODO: Insert CheckID method here
		err = h.service.Store.CheckID(id)
		if err != nil {
			h.SendErr(w, err, http.StatusBadRequest)
			return
		}
		// Insert method Delete
		err = h.service.Store.Delete(id)
		if err != nil {
			h.SendErr(w, err, http.StatusInternalServerError)
			return
		}
		h.logger.Infof("sent response via handler Task (method %s)", r.Method)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprint(w, "{}")
	default:
		err := fmt.Errorf("no request method")
		h.SendErr(w, err, http.StatusInternalServerError)
		return
	}
}