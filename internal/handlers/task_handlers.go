package handlers

import (
	"async-task-service/internal/task"
	"encoding/json"
	"net/http"
	"strings"
)

func RegisterTaskHandlers(mux *http.ServeMux, manager *task.Manager) {
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateTask(w, r, manager)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/tasks/")
		switch r.Method {
		case http.MethodGet:
			handleGetTask(w, r, manager, id)
		case http.MethodDelete:
			handleDeleteTask(w, r, manager, id)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func handleCreateTask(w http.ResponseWriter, r *http.Request, manager *task.Manager) {
	t := manager.CreateTask()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": t.ID})
}

func handleGetTask(w http.ResponseWriter, r *http.Request, manager *task.Manager, id string) {
	t, err := manager.Get(id)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	dto := t.ToDTO()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto)
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request, manager *task.Manager, id string) {
	err := manager.Delete(id)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
