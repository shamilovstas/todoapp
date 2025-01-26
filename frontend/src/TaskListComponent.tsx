import './TaskList.css'
import React, {ChangeEvent, useEffect, useState} from "react";
import {NavLink} from "react-router";
import {TaskList, apiUrl} from "./types.ts";

interface InputProps {
    input: string;
    handleChange: (e: ChangeEvent<HTMLInputElement>) => void;
    onKeyUp: (event: React.KeyboardEvent<HTMLInputElement>) => void;
}

function TaskListInput({input, handleChange, onKeyUp}: InputProps) {

    const placeholder = "Enter task"
    return <input
        className="tasklist-input input-light"
        placeholder={placeholder}
        // onFocus={(e) => e.target.placeholder = ""}
        // onBlur={(e) => e.target.placeholder = placeholder}
        type="text"
        onChange={handleChange}
        onKeyUp={onKeyUp}
        value={input}/>
}

function List({list}: { list: TaskList[] }) {
    return list.map((item) => <li key={item.id}>
        <NavLink
            className={({isActive}) => ["task-list-link", isActive ? "active" : ""].join(" ")}
            to={`/tasklists/${item.id}`}>{item.name}
        </NavLink>
    </li>)
}

function TaskListComponent() {

    const [list, setList] = useState<TaskList[]>([]);
    const [input, setInput] = useState<string>('')

    useEffect(() => {
        fetch(`${apiUrl}/tasklists`)
            .then(res => {

                const json = res.json()
                return json;
            })
            .then(json => json.taskLists)
            .then(taskList => setList(taskList))
            .catch(err => console.log(err))
    }, [])

    const handleInput = (e: ChangeEvent<HTMLInputElement>) => {
        setInput(e.target.value);
    }

    const addTaskList = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key !== 'Enter') {
            return
        }
        if (input.trim() !== '') {
            const data = {"name": input}
            fetch(`${apiUrl}/tasklists`, {
                method: "POST",
                body: JSON.stringify(data),
                headers: {"Content-Type": "application/json"},
            })
                .then(res => res.json())
                .then(taskList => {
                    setList([...list, taskList]);
                    setInput("")
                })
                .catch(error => console.log(error))
        }
    }

    return <div className="tasklist-container">
        <h1>My lists</h1>
        <ul className="task-list"><List list={list}/></ul>
        <div className="plus-input">
            <TaskListInput input={input} handleChange={handleInput} onKeyUp={addTaskList}/>
        </div>
    </div>
}

export default TaskListComponent