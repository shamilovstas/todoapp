package tasks

type Task struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	IsCompleted bool   `json:"completed"`
}

type TaskList struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}

func NewTask(id int, name string, completed bool) *Task {
	return &Task{Id: id, Name: name, IsCompleted: completed}
}

func NewTaskList(id int, name string) *TaskList {
	return &TaskList{Id: id, Name: name}
}
