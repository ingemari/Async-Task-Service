package main

import (
	"log"
	"net/http"

	"async-task-service/internal/handlers"
	"async-task-service/internal/task"
)

func main() {
	manager := task.NewManager()

	mux := http.NewServeMux()
	handlers.RegisterTaskHandlers(mux, manager)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
