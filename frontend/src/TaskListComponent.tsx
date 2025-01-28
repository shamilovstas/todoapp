import './TaskList.css'
import React, {ChangeEvent, useEffect, useState} from "react";
import {apiUrl, TaskList} from "./types.ts";
import TasksComponent from "./TasksComponent.tsx";

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

interface ListProps {
    list: TaskList[];
    selectedId: number | undefined;
    onClick: (t: TaskList) => void;
}

function List({list, selectedId, onClick}: ListProps) {
    return list.map(
        (item) =>
            <li
                key={item.id}
                className={item.id === selectedId ? "active" : ""}
                onClick={() => {
                    onClick(item)
                }}>
                {item.name}
            </li>
    )
}

function TaskListComponent() {

    const [list, setList] = useState<TaskList[]>([]);
    const [selectedList, setSelectedList] = useState<TaskList | null>();
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

    const deleteTaskList = (taskList: TaskList) => {
        fetch(`${apiUrl}/tasklists/${taskList.id}`, {
            method: "DELETE",
            headers: {"Content-Type": "application/json"}
        })
            .then(() => {
                setList(list.filter((task) => task.id !== taskList.id))
                setSelectedList(null)
            })
            .catch(err => console.log(err))
    }

    return <div className="tasklist-container">
        <div className="list-container">
            <h1>My lists</h1>
            <ul className="task-list"><List list={list} selectedId={selectedList?.id}
                                            onClick={(t: TaskList) => setSelectedList(t)}/></ul>
            <div className="plus-input">
                <TaskListInput input={input} handleChange={handleInput} onKeyUp={addTaskList}/>
            </div>
        </div>
        <div className="details-container">
            {selectedList ?
                <>
                    <TasksComponent id={selectedList.id}/>
                    <div className="buttons">
                        <button type="button">Clear completed tasks</button>
                        <button type="button" onClick={() => deleteTaskList(selectedList)}>Delete list</button>
                    </div>
                </>
                : null
            }

        </div>
    </div>
}

export default TaskListComponent