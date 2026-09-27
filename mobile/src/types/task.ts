export type TaskStatus = 'todo' | 'in_progress' | 'done';

export interface Task {
  id: number;
  title: string;
  description: string;
  status: TaskStatus;
  assignee: string;
  created_at: string;
  updated_at: string;
}

export interface TaskListResponse {
  items: Task[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface TaskFilters {
  page: number;
  limit: number;
  keyword: string;
  status: '' | TaskStatus;
  assignee: string;
  sort: string;
}

export interface UpdateTaskPayload {
  title: string;
  description: string;
  status: TaskStatus;
  assignee: string;
}
