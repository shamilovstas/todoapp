package tasks

type Task struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TaskList struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}

func NewTask(id int, name string) *Task {
	return &Task{Id: id, Name: name}
}

func NewTaskList(id int, name string) *TaskList {
	return &TaskList{Id: id, Name: name}
}
