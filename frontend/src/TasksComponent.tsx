import './Tasks.css'
import React, {ChangeEvent, useEffect, useState} from "react";
import {useParams} from "react-router";
import {Task, TaskList, apiUrl} from "./types.ts";

interface TaskListProps {
    taskList: TaskList,
    handleToggle: (task: Task, isCompleted: boolean) => void
}

function RenderedTasks({taskList, handleToggle}: TaskListProps) {

    return taskList.tasks.map((item: Task) => {
        return <li className="task-item" key={item.id}>
            <input
                type="checkbox"
                id={`checkbox-${item.id}`}
                checked={item.completed}
                onChange={(e) => handleToggle(item, e.target.checked)}
            />
            <label htmlFor={`checkbox-${item.id}`}>{item.name}</label>
        </li>
    });
}

interface TaskInputProps {
    input: string;
    handleChange: (e: ChangeEvent<HTMLInputElement>) => void;
    onKeyUp: (event: React.KeyboardEvent<HTMLInputElement>) => void;
}

function NewTaskInput({input, handleChange, onKeyUp}: TaskInputProps) {
    return <input
        className="input-dark"
        onChange={handleChange}
        onKeyUp={onKeyUp}
        value={input}/>
}

function TasksComponent() {

    const {id} = useParams()
    const [taskList, setTaskList] = useState<TaskList>({id: -1, name: "", tasks: [], remaining: 0});
    const [input, setInput] = useState<string>('')

    const handleInput = (e: ChangeEvent<HTMLInputElement>) => {
        setInput(e.target.value)
    }

    const saveTask = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key !== 'Enter') {
            return;
        }
        if (input.trim() !== '') {
            const data = {"name": input};
            fetch(`${apiUrl}/tasklists/${id}/tasks`, {
                method: "POST",
                body: JSON.stringify(data),
                headers: {"Content-Type": "application/json"},
            })
                .then(res => res.json())
                .then(task => {
                    taskList.tasks.push(task)
                    setTaskList(taskList)
                    setInput("")
                })
                .catch(error => console.log(error))
        }
    }

    useEffect(() => {
        fetch(`${apiUrl}/tasklists/${id}`)
            .then(res => res.json())
            .then(data => {
                setTaskList(data)
            })
            .catch(error => console.log(error))
    }, [id])

    const handleToggle = (oldTask: Task, isCompleted: boolean) => {
        oldTask.completed = isCompleted;

        fetch(apiUrl + "/tasks", {
            method: "PUT",
            body: JSON.stringify(oldTask),
            headers: {"Content-Type": "application/json"}
        })
            .then(() => {
                const updatedTasks = taskList.tasks.map(task => {
                    if (task.id === task.id) {
                        return {...task, completed: isCompleted};
                    } else {
                        return task;
                    }
                })


                let remaining = taskList.remaining;
                if (isCompleted) {
                    remaining --;
                } else {
                    remaining ++;
                }
                const newTaskList = {...taskList, remaining, updatedTasks}
                setTaskList(newTaskList)
            })
            .catch(e => console.log(e))

    }

    return <div className="task-container">
        <div className="task-header">
            <h2>{taskList.name}</h2>
            <p>{taskList.remaining} tasks remaining</p>
        </div>

        <div className="task-content">
            <ul className="task-items">
                <RenderedTasks taskList={taskList} handleToggle={handleToggle}/>
            </ul>

            <div className="plus-input">
                <NewTaskInput onKeyUp={saveTask} handleChange={handleInput} input={input}/>
            </div>
        </div>


    </div>
}

export default TasksComponent;