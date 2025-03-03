export interface Task {
    id: number;
    name: string;
    completed: boolean;
}

export interface TaskList {
    id: number;
    name: string;
    tasks: Task[];
    remaining: number;
}

// export const apiUrl = "http://localhost:29292";
export const apiUrl = "/api";