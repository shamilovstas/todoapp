import './App.css'
import TaskListComponent from "./TaskListComponent.tsx";
import {MemoryRouter} from "react-router";

function App() {
    return <div className="container">
        <MemoryRouter>
            <h1 className="header">Stuff I want to do</h1>
            <TaskListComponent/>

        </MemoryRouter>
    </div>
}

export default App
