import { apiClient } from '$lib/api/client';
import type { Environment, CreateEnvironmentRequest, EnvironmentsResponse } from '../types';

export const environmentsApi = {
  async list(projectId: string): Promise<Environment[]> {
    const response = await apiClient.get<EnvironmentsResponse>(
      `/projects/${projectId}/environments`
    );
    return response.environments;
  },

  async get(projectId: string, environmentId: string): Promise<Environment> {
    return apiClient.get<Environment>(`/projects/${projectId}/environments/${environmentId}`);
  },

  async create(projectId: string, data: CreateEnvironmentRequest): Promise<Environment> {
    return apiClient.post<Environment>(`/projects/${projectId}/environments`, data);
  },

  async update(
    projectId: string,
    environmentId: string,
    data: Partial<CreateEnvironmentRequest>
  ): Promise<Environment> {
    return apiClient.put<Environment>(
      `/projects/${projectId}/environments/${environmentId}`,
      data
    );
  },

  async delete(projectId: string, environmentId: string): Promise<void> {
    return apiClient.delete<void>(`/projects/${projectId}/environments/${environmentId}`);
  }
};
