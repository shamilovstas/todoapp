import './Tasks.css'
import React, {ChangeEvent, useEffect, useState} from "react";
import {useParams} from "react-router";
import {Task, TaskList, apiUrl} from "./types.ts";

function RenderedTasks({list}: { list: Task[] }) {
    return list.map((item: Task) => {
        return <li className="task-item" key={item.id}>{item.name}</li>
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
    const [taskList, setTaskList] = useState<TaskList>({id: -1, name: "", tasks: []});
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
    return <div className="task-container">
        <div className="task-header">
            <h2>{taskList.name}</h2>
            <p>{taskList.tasks.length} tasks remaining</p>
        </div>

        <div className="task-content">
            <ul className="task-list">
                <RenderedTasks list={taskList.tasks}/>
            </ul>
        </div>

        <div className="plus-input">
            <NewTaskInput onKeyUp={saveTask} handleChange={handleInput} input={input}/>
        </div>
    </div>
}

export default TasksComponent;