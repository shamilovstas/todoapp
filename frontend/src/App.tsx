import './App.css'
import TaskListComponent from "./TaskListComponent.tsx";
import TasksComponent from "./TasksComponent.tsx";
import {MemoryRouter, Navigate, Route, Routes} from "react-router";

function App() {
    return <div className="container">
        <MemoryRouter>
            <h1 className="header">Stuff I want to do</h1>
            <TaskListComponent/>
            <Routes>
                <Route path="/" element={<Navigate to={`/tasklists/1`} replace /> }/>
                <Route path="/tasklists/:id" element={<TasksComponent/>}/>
            </Routes>
        </MemoryRouter>
    </div>
}

export default App
