package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"todo-api/api/resource/tasks"
	"todo-api/database"
	"todo-api/middleware"
)

//type TaskListResponse struct {
//	TaskLists []tasks.TaskList `json:"taskLists"`
//}
//
//type TasksResponse struct {
//	Id             int          `json:"id"`
//	Name           string       `json:"name"`
//	Tasks          []tasks.Task `json:"tasks"`
//	RemainingTasks int          `json:"remaining"`
//}

func taskHandler(api *tasks.Api) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			api.DeleteTask(w, r)
		case http.MethodPut:
			api.UpdateTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func tasksHandler(api *tasks.Api) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			api.UpdateTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func taskListHandler(api *tasks.Api) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.GetTaskListById(w, r)
		case http.MethodDelete:
			api.DeleteTaskList(w, r)
		case http.MethodPost:
			api.CreateTaskList(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func taskListsHandler(api *tasks.Api) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.GetTaskLists(w, r)
		case http.MethodPost:
			api.CreateTaskList(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func ignoreCors(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")                                // Allow all origins (*). Replace with specific origin in production.
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // Allowed methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")     // Allowed headers

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass to the next handler
		n.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	pool, err := database.CreatePool()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create database connection pool: %v\n", err)
		os.Exit(1)
	}

	if err := pool.Ping(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	api := tasks.New(pool)

	mux.HandleFunc("/tasklists", taskListsHandler(api))
	mux.HandleFunc("/tasklists/{id}/tasks", api.CreateTask)
	mux.HandleFunc("/tasklists/{id}", taskListHandler(api))
	mux.HandleFunc("/tasks/{id}", taskHandler(api))
	mux.HandleFunc("/tasklists/{id}/tasks/completed", api.DeleteCompletedTasks)
	if err := http.ListenAndServe(":8080", middleware.AddMiddleware(mux)); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't start server: %v\n", err)
		os.Exit(1)
	}
}
