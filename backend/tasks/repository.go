package tasks

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
	if err := row.Scan(&name); err != nil {
		return TaskList{}, err
	}
	return *NewTaskList(id, name), nil
}

func (repo TaskRepository) IsTaskListExists(id int) (bool, error) {
	query := `SELECT id FROM tasklists WHERE id=@taskListId`
	args := pgx.NamedArgs{"taskListId": id}
	rows, err := repo.db.Query(context.Background(), query, args)
	defer rows.Close()
	if err != nil {
		return false, err
	}

	rowsCount := 0
	for rows.Next() {
		rowsCount++
	}

	return rowsCount != 0, nil
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
		var id int
		var name string
		var completed bool

		if err := rows.Scan(&id, &name, &completed); err != nil {
			return nil, err
		}

		task := NewTask(id, name, completed)
		taskArray = append(taskArray, *task)
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

func (repo TaskRepository) UpdateTask(ctx context.Context, task Task) (int, error) {
	query := `UPDATE tasks SET name=@taskName, completed=@isCompleted WHERE tasks.id=@taskId RETURNING id`
	args := pgx.NamedArgs{
		"taskName":    task.Name,
		"isCompleted": task.IsCompleted,
		"taskId":      task.Id,
	}

	row := repo.db.QueryRow(ctx, query, args)
	insertedId := -1

	if err := row.Scan(&insertedId); err != nil {
		return -1, err
	}

	return insertedId, nil
}
