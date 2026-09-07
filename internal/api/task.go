package api

import (
	"encoding/json"
	"net/http"
	"scheduler/internal/db"
	"strconv"
	"time"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSONError(w, "title is required", http.StatusBadRequest)
		return
	}

	now := time.Now()
	nowStr := now.Format(dateFormat)
	nowDate, _ := time.Parse(dateFormat, nowStr)

	if task.Date == "" {
		task.Date = nowStr
	}
	parsedDate, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		writeJSONError(w, "invalid date format, expected YYYYMMDD", http.StatusBadRequest)
		return
	}

	var nextDate string
	if task.Repeat != "" {
		next, err := NextDate(nowStr, task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, "invalid repeat rule: "+err.Error(), http.StatusBadRequest)
			return
		}
		nextDate = next
	}

	if parsedDate.Before(nowDate) {
		if task.Repeat == "" {
			task.Date = nowStr
		} else {
			task.Date = nextDate
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONError(w, "database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"id": id})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, "id is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, "invalid id", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, "task not found", http.StatusNotFound)
		return
	}
	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "invalid JSON format", http.StatusBadRequest)
		return
	}
	if task.ID == 0 {
		writeJSONError(w, "id is required", http.StatusBadRequest)
		return
	}

	if _, err := db.GetTask(task.ID); err != nil {
		writeJSONError(w, "task not found", http.StatusNotFound)
		return
	}

	if task.Title == "" {
		writeJSONError(w, "title is required", http.StatusBadRequest)
		return
	}

	now := time.Now()
	nowStr := now.Format(dateFormat)
	nowDate, _ := time.Parse(dateFormat, nowStr)

	if task.Date == "" {
		task.Date = nowStr
	}
	parsedDate, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		writeJSONError(w, "invalid date format, expected YYYYMMDD", http.StatusBadRequest)
		return
	}

	var nextDate string
	if task.Repeat != "" {
		next, err := NextDate(nowStr, task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, "invalid repeat rule: "+err.Error(), http.StatusBadRequest)
			return
		}
		nextDate = next
	}

	if parsedDate.Before(nowDate) {
		if task.Repeat == "" {
			task.Date = nowStr
		} else {
			task.Date = nextDate
		}
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJSONError(w, "database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{})
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error":"failed to encode JSON"}`, http.StatusInternalServerError)
	}
}

func writeJSONError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSONError(w, "database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"tasks": tasks})
}
