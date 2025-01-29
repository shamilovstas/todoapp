package tasks

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTaskListNotFound = errors.New("task list not found")
var ErrTaskNotFound = errors.New("task not found")

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (repo TaskRepository) GetTaskLists(ctx context.Context) ([]TaskList, error) {
	query := `SELECT id, name FROM tasklists`
	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	taskLists := make([]TaskList, 0)
	for rows.Next() {
		var id int
		var name string

		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}

		taskList := NewTaskList(id, name)
		taskLists = append(taskLists, *taskList)
	}
	return taskLists, nil
}

func (repo TaskRepository) InsertTaskList(ctx context.Context, tl TaskList) (int, error) {
	query := `INSERT INTO tasklists (name) VALUES (@taskName) RETURNING id`
	args := pgx.NamedArgs{
		"taskName": tl.Name,
	}

	row := repo.db.QueryRow(ctx, query, args)
	insertedId := -1

	if err := row.Scan(&insertedId); err != nil {
		return -1, err
	}

	return insertedId, nil
}

func (repo TaskRepository) GetTaskListById(id int) (TaskList, error) {
	query := `SELECT name FROM tasklists WHERE id=@taskListId`
	args := pgx.NamedArgs{"taskListId": id}
	row := repo.db.QueryRow(context.Background(), query, args)
	var name string
	var taskListError error
	var taskList TaskList
	if err := row.Scan(&name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			taskListError = ErrTaskListNotFound
		} else {
			taskListError = err
		}
	} else {
		taskList = *NewTaskList(id, name)
	}
	return taskList, taskListError
}

func (repo TaskRepository) GetTasks(ctx context.Context, listId int) ([]Task, error) {
	args := pgx.NamedArgs{"listId": listId}
	query := `SELECT id, name, completed FROM tasks WHERE tasks.listId=@listId`
	rows, err := repo.db.Query(ctx, query, args)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	taskArray := make([]Task, 0)
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Name, &task.IsCompleted); err != nil {
			return nil, err
		}
		taskArray = append(taskArray, task)
	}
	return taskArray, nil
}

func (repo TaskRepository) InsertTask(ctx context.Context, task Task, listId int) (int, error) {
	query := `INSERT INTO tasks (name, listId, completed) VALUES (@taskName, @listId, @isCompleted) RETURNING id`
	args := pgx.NamedArgs{
		"taskName":    task.Name,
		"listId":      listId,
		"isCompleted": task.IsCompleted,
	}

	row := repo.db.QueryRow(ctx, query, args)
	insertedId := -1

	if err := row.Scan(&insertedId); err != nil {
		return -1, err
	}

	return insertedId, nil
}

func (repo TaskRepository) UpdateTask(ctx context.Context, taskId int, task Task) (int, error) {
	query := `UPDATE tasks SET name=@taskName, completed=@isCompleted WHERE tasks.id=@taskId RETURNING id`
	args := pgx.NamedArgs{
		"taskName":    task.Name,
		"isCompleted": task.IsCompleted,
		"taskId":      taskId,
	}

	row := repo.db.QueryRow(ctx, query, args)
	insertedId := -1
	var updateError error

	if err := row.Scan(&insertedId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			updateError = ErrTaskNotFound
		} else {
			updateError = err
		}
	}
	return insertedId, updateError
}

func (repo TaskRepository) DeleteTask(ctx context.Context, taskId int) error {
	query := `DELETE FROM tasks WHERE tasks.id=@taskId`
	args := pgx.NamedArgs{
		"taskId": taskId,
	}

	_, err := repo.db.Exec(ctx, query, args)
	return err
}

func (repo TaskRepository) DeleteTaskList(ctx context.Context, taskListId int) error {
	query := `DELETE FROM taskLists WHERE taskLists.id=@taskListId`
	args := pgx.NamedArgs{
		"taskListId": taskListId,
	}
	_, err := repo.db.Exec(ctx, query, args)
	return err
}

func (repo TaskRepository) DeleteCompletedTasks(ctx context.Context, taskListId int) error {

	args := pgx.NamedArgs{
		"taskListId": taskListId,
	}
	query := `DELETE FROM tasks WHERE tasks.listId=@taskListId AND tasks.completed IS TRUE`

	_, err := repo.db.Exec(ctx, query, args)
	return err
}
