import './Tasks.css'
import React, {ChangeEvent, useEffect, useMemo, useState} from "react";
import {Task, TaskList, apiUrl} from "./types.ts";

interface TaskListProps {
    tasks: Task[],
    handleToggle: (task: Task, isCompleted: boolean) => void
}

function RenderedTasks({tasks, handleToggle}: TaskListProps) {

    return tasks.map((item: Task) => {
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

interface TasksComponentProps {
    taskList: TaskList
}

function TasksComponent({taskList}: TasksComponentProps) {

    const [tasks, setTasks] = useState<Task[]>([])
    const [input, setInput] = useState<string>('')

    const remaining = useMemo(() => tasks.filter((task) => !task.completed).length, [tasks])

    useEffect(() => {
        fetch(`${apiUrl}/tasklists/${taskList.id}`)
            .then(res => res.json())
            .then(data => {
                setTasks(data.tasks)
            })
            .catch(error => console.log(error))
    }, [taskList])

    const handleInput = (e: ChangeEvent<HTMLInputElement>) => {
        setInput(e.target.value)
    }

    const saveTask = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key !== 'Enter') {
            return;
        }
        if (input.trim() !== '') {
            const data = {"name": input};
            fetch(`${apiUrl}/tasklists/${taskList.id}/tasks`, {
                method: "POST",
                body: JSON.stringify(data),
                headers: {"Content-Type": "application/json"},
            })
                .then(res => res.json())
                .then(task => {
                    setTasks([...tasks, task]);
                    setInput("")
                })
                .catch(error => console.log(error))
        }

    }

    const handleToggle = (oldTask: Task, isCompleted: boolean) => {
        oldTask.completed = isCompleted;

        fetch(apiUrl + `/tasks/${oldTask.id}`, {
            method: "PUT",
            body: JSON.stringify(oldTask),
            headers: {"Content-Type": "application/json"}
        })
            .then(() => {
                const updatedTasks = tasks.map(task => {
                    if (oldTask.id === task.id) {
                        return {...task, completed: isCompleted};
                    } else {
                        return task;
                    }
                })
                setTasks(updatedTasks)
            })
            .catch(e => console.log(e))

    }

    return <div className="task-container">
        <div className="task-header">
            <h2>{taskList.name}</h2>
            <p>{remaining} tasks remaining</p>
        </div>

        <div className="task-content">
            <ul className="task-items">
                <RenderedTasks tasks={tasks} handleToggle={handleToggle}/>
            </ul>

            <div className="plus-input">
                <NewTaskInput onKeyUp={saveTask} handleChange={handleInput} input={input}/>
            </div>
        </div>
    </div>
}

export default TasksComponent;