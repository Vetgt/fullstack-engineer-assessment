import { Task, TaskFilters, TaskListResponse, UpdateTaskPayload } from '../types/task';

const BASE_URL = (process.env.EXPO_PUBLIC_API_URL ?? 'http://localhost:8080/api').replace(/\/$/, '');

class ApiError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.status = status;
    this.name = 'ApiError';
  }
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const response = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers ?? {}),
    },
  });

  if (!response.ok) {
    let message = `Request failed with HTTP ${response.status}`;
    try {
      const body = (await response.json()) as { error?: { message?: string } };
      message = body.error?.message ?? message;
    } catch {
      // Keep the default message when the response is not JSON.
    }
    throw new ApiError(message, response.status);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return (await response.json()) as T;
}

export async function getTasks(filters: TaskFilters): Promise<TaskListResponse> {
  const params = new URLSearchParams({
    page: String(filters.page),
    limit: String(filters.limit),
    sort: filters.sort,
  });

  if (filters.keyword.trim()) {
    params.set('keyword', filters.keyword.trim());
  }
  if (filters.status) {
    params.set('status', filters.status);
  }
  if (filters.assignee.trim()) {
    params.set('assignee', filters.assignee.trim());
  }

  return request<TaskListResponse>(`/tasks?${params.toString()}`);
}

export async function updateTask(id: number, payload: UpdateTaskPayload): Promise<Task> {
  return request<Task>(`/tasks/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
}

export { ApiError };
