package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"todo-api/database"
	"todo-api/tasks"
)

var taskRepository tasks.TaskRepository

type TaskListResponse struct {
	TaskLists []tasks.TaskList `json:"taskLists"`
}

func tasklistsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		taskLists, err := taskRepository.GetTaskLists(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := TaskListResponse{TaskLists: taskLists}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		var tl tasks.TaskList

		err := json.NewDecoder(r.Body).Decode(&tl)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := taskRepository.InsertTaskList(context.Background(), tl)
		tl.Id = id

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(tl); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getTaskList(w http.ResponseWriter, listId int) error {
	taskList, err := taskRepository.GetTaskListById(listId)
	if err != nil {
		return NewHttpError(http.StatusNotFound, fmt.Errorf("task list not found: %w", err).Error())
	}
	taskArray, err := taskRepository.GetTasks(context.Background(), listId)
	if err != nil {
		return err
	}

	taskList.Tasks = taskArray

	b, _ := json.Marshal(taskList)
	log.Println(string(b))

	if err := json.NewEncoder(w).Encode(taskList); err != nil {
		return err
	}
	return nil
}

func postTask(w http.ResponseWriter, r *http.Request, listId int) error {
	var task tasks.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		return err
	}

	id, err := taskRepository.InsertTask(context.Background(), task, listId)
	if err != nil {
		return err
	}
	task.Id = id
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		return err
	}
	return nil
}

type HttpError struct {
	Code    int
	Message string
}

func NewHttpError(code int, message string) *HttpError {
	return &HttpError{Code: code, Message: message}
}

func (error HttpError) Error() string {
	return error.Message
}

func checkTaskListExists(id int) error {
	taskListExists, err := taskRepository.IsTaskListExists(id)
	if err != nil {
		return NewHttpError(http.StatusInternalServerError, err.Error())
	}
	if !taskListExists {
		return NewHttpError(http.StatusNotFound, "Task list not found")
	}
	return nil
}

func getTaskListHandler(w http.ResponseWriter, r *http.Request) {
	listId, idAtoiErr := strconv.Atoi(r.PathValue("id"))
	if idAtoiErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := checkTaskListExists(listId); err != nil {
		log.Printf("error: %v\n", err)
		var herr *HttpError
		ok := errors.As(err, &herr)
		if ok {
			http.Error(w, herr.Message, herr.Code)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	if err := getTaskList(w, listId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	listId, idAtoiErr := strconv.Atoi(r.PathValue("id"))
	if idAtoiErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := checkTaskListExists(listId); err != nil {
		fmt.Printf("error: %v\n", err)
		var herr *HttpError
		ok := errors.As(err, &herr)
		if ok {
			http.Error(w, herr.Message, herr.Code)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	if err := postTask(w, r, listId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := updateTask(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) error {
	var task tasks.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		return err
	}

	id, err := taskRepository.UpdateTask(context.Background(), task)
	if err != nil {
		return err
	}
	task.Id = id
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		return err
	}
	return nil
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

	taskRepository = *tasks.NewTaskRepository(pool)

	mux.HandleFunc("/tasklists", tasklistsHandler)
	mux.HandleFunc("/tasklists/{id}/tasks", createTaskHandler)
	mux.HandleFunc("/tasklists/{id}", getTaskListHandler)
	mux.HandleFunc("/tasks", updateTaskHandler)
	if err := http.ListenAndServe(":8080", ignoreCors(mux)); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't start server: %v\n", err)
		os.Exit(1)
	}
}
