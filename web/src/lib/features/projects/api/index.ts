import { apiClient } from '$lib/api/client';
import type { Project, CreateProjectRequest, ProjectsResponse } from '../types';

export const projectsApi = {
  async list(): Promise<Project[]> {
    const response = await apiClient.get<ProjectsResponse>('/projects');
    return response.projects;
  },

  async get(id: string): Promise<Project> {
    return apiClient.get<Project>(`/projects/${id}`);
  },

  async create(data: CreateProjectRequest): Promise<Project> {
    return apiClient.post<Project>('/projects', data);
  },

  async update(id: string, data: Partial<CreateProjectRequest>): Promise<Project> {
    return apiClient.put<Project>(`/projects/${id}`, data);
  },

  async delete(id: string): Promise<void> {
    return apiClient.delete<void>(`/projects/${id}`);
  },
};
