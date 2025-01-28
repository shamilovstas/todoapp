package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"todo-api/apperrors"
	"todo-api/database"
	"todo-api/tasks"
)

var taskRepository tasks.TaskRepository

type TaskListResponse struct {
	TaskLists []tasks.TaskList `json:"taskLists"`
}

type TasksResponse struct {
	Id             int          `json:"id"`
	Name           string       `json:"name"`
	Tasks          []tasks.Task `json:"tasks"`
	RemainingTasks int          `json:"remaining"`
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
		return apperrors.NewHttpError(http.StatusNotFound, fmt.Errorf("task list not found: %w", err))
	}
	taskArray, err := taskRepository.GetTasks(context.Background(), listId)
	if err != nil {
		return err
	}
	taskList.Tasks = taskArray
	response := TasksResponse{
		Id:             taskList.Id,
		Name:           taskList.Name,
		Tasks:          taskList.Tasks,
		RemainingTasks: taskList.GetRemainingTasksCount(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
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

func checkTaskListExists(id int) error {
	taskListExists, err := taskRepository.IsTaskListExists(id)
	if err != nil {
		return apperrors.NewHttpError(http.StatusInternalServerError, err)
	}
	if !taskListExists {
		return apperrors.NewHttpError(http.StatusNotFound, errors.New("task list not found"))
	}
	return nil
}

func deleteTaskList(w http.ResponseWriter, taskListId int) error {
	return taskRepository.DeleteTaskList(context.Background(), taskListId)
}

func taskListHandler(w http.ResponseWriter, r *http.Request) {
	listId, idAtoiErr := strconv.Atoi(r.PathValue("id"))
	if idAtoiErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := checkTaskListExists(listId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

	}

	switch r.Method {
	case "GET":
		if err := getTaskList(w, listId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "DELETE":
		if err := deleteTaskList(w, listId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := postTask(w, r, listId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		if err := updateTask(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "DELETE":
		if err := deleteTask(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) error {
	taskId, idAtoiErr := strconv.Atoi(r.PathValue("id"))
	if idAtoiErr != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, errors.New("task not found"))
	}
	var task tasks.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		return err
	}

	id, err := taskRepository.UpdateTask(context.Background(), taskId, task)
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

func deleteTask(w http.ResponseWriter, r *http.Request) error {
	taskId, idAtoiErr := strconv.Atoi(r.PathValue("id"))
	if idAtoiErr != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, errors.New("task not found"))
	}

	err := taskRepository.DeleteTask(context.Background(), taskId)
	if err == nil {
		w.WriteHeader(http.StatusOK)
	}
	return err
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

func completedTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	taskListId, idAtoiErr := strconv.Atoi(r.PathValue("id"))
	if idAtoiErr != nil {
		http.Error(w, "task list not found", http.StatusNotFound)
		return
	}

	if err := taskRepository.DeleteCompletedTasks(context.Background(), taskListId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
	mux.HandleFunc("/tasklists/{id}", taskListHandler)
	mux.HandleFunc("/tasks/{id}", taskHandler)
	mux.HandleFunc("/tasklists/{id}/tasks/completed", completedTasksHandler)
	if err := http.ListenAndServe(":8080", ignoreCors(mux)); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't start server: %v\n", err)
		os.Exit(1)
	}
}
