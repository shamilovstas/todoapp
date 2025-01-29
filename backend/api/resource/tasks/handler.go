package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"strconv"
	e "todo-api/api/common/err"
	"todo-api/tasks"
)

type Api struct {
	repository TaskRepository
}

func New(db *pgxpool.Pool) *Api {
	return &Api{repository: NewTaskRepository(db)}
}

func (api *Api) CreateTask(w http.ResponseWriter, r *http.Request) {
	listId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		e.BadRequest(w, errors.New("id was not an integer"))
		return
	}

	var task tasks.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		e.ServerError(w, e.RespErrJsonDecode)
		return
	}

	taskId, err := api.repository.InsertTask(context.Background(), task, listId)
	if err != nil {
		if errors.Is(err, ErrTaskListNotFound) {
			e.NotFound(w, err)
		} else {
			e.ServerError(w, e.RespErrDbInsert)
		}
		return
	}

	task.Id = taskId
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		e.ServerError(w, e.RespErrJsonEncode)
	}
}

func (api *Api) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		e.BadRequest(w, errors.New("id was not an integer"))
		return
	}
	var task tasks.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		e.ServerError(w, e.RespErrJsonDecode)
		return
	}

	id, err := api.repository.UpdateTask(context.Background(), taskId, task)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			e.NotFound(w, err)
		} else {
			e.ServerError(w, e.RespErrDbUpdate)
		}
		return
	}
	task.Id = id
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		e.ServerError(w, e.RespErrJsonEncode)
	}
}

func (api *Api) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		e.BadRequest(w, errors.New("id was not an integer"))
		return
	}

	if err := api.repository.DeleteTask(context.Background(), taskId); err != nil {
		e.ServerError(w, e.RespErrDbDelete)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (api *Api) GetTaskListById(w http.ResponseWriter, r *http.Request) {
	listId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		e.BadRequest(w, errors.New("id was not an integer"))
		return
	}

	taskList, err := api.repository.GetTaskListById(context.Background(), listId)
	if err != nil {
		if errors.Is(err, ErrTaskListNotFound) {
			e.NotFound(w, err)
		} else {
			e.ServerError(w, e.RespErrDbAccess)
		}
		return
	}
	taskArray, err := api.repository.GetTasks(context.Background(), listId)
	if err != nil {
		e.ServerError(w, e.RespErrDbAccess)
		return
	}
	taskList.Tasks = taskArray

	if err := json.NewEncoder(w).Encode(taskList); err != nil {
		e.ServerError(w, e.RespErrJsonEncode)
	}
}

func (api *Api) CreateTaskList(w http.ResponseWriter, r *http.Request) {
	var tl tasks.TaskList

	err := json.NewDecoder(r.Body).Decode(&tl)

	if err != nil {
		e.BadRequest(w, e.RespErrJsonDecode)
		return
	}

	id, err := api.repository.InsertTaskList(context.Background(), tl)
	if err != nil {
		e.ServerError(w, e.RespErrDbInsert)
		return
	}
	tl.Id = id

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(tl); err != nil {
		e.ServerError(w, e.RespErrJsonEncode)
	}
}

func (api *Api) DeleteTaskList(w http.ResponseWriter, r *http.Request) {
	listId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		e.BadRequest(w, errors.New("id was not an integer"))
		return
	}

	if err := api.repository.DeleteTaskList(context.Background(), listId); err != nil {
		e.ServerError(w, e.RespErrDbDelete)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (api *Api) GetTaskLists(w http.ResponseWriter, r *http.Request) {
	taskLists, err := api.repository.GetTaskLists(context.Background())
	if err != nil {
		e.ServerError(w, e.RespErrDbAccess)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(taskLists); err != nil {
		e.ServerError(w, e.RespErrJsonEncode)
	}
}

func (api *Api) DeleteCompletedTasks(w http.ResponseWriter, r *http.Request) {
	listId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		e.BadRequest(w, errors.New("id was not an integer"))
		return
	}

	if err := api.repository.DeleteCompletedTasks(context.Background(), listId); err != nil {
		e.ServerError(w, e.RespErrDbDelete)
		return
	}

	api.GetTaskListById(w, r)
}
