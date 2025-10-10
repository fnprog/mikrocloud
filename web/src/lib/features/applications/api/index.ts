import { apiClient } from '$lib/api/client';
import type { Application, ApplicationsResponse, CreateApplicationRequest, DeleteApplicationRequest, SwitchApplicationStateRequest, UpdateGeneralSettingsRequest } from '../types';


export const applicationsApi = {
  async list(projectId: string): Promise<Application[]> {
    const response = await apiClient.get<ApplicationsResponse>(
      `/projects/${projectId}/applications`
    );
    return response.applications;
  },

  async get(projectId: string, applicationId: string): Promise<Application> {
    return apiClient.get<Application>(`/projects/${projectId}/applications/${applicationId}`);
  },

  async create(data: CreateApplicationRequest): Promise<Application> {
    return apiClient.post<Application>(`/projects/${data.project_id}/applications`, data);
  },

  async delete(data: DeleteApplicationRequest): Promise<void> {
    return apiClient.delete<void>(`/projects/${data.project_id}/applications/${data.id}`);
  },

  async start(data: SwitchApplicationStateRequest): Promise<void> {
    return apiClient.post<void>(`/projects/${data.project_id}/applications/${data.application_id}/start`, {});
  },

  async stop(data: SwitchApplicationStateRequest): Promise<void> {
    return apiClient.post<void>(`/projects/${data.project_id}/applications/${data.application_id}/stop`, {});
  },

  async restart(data: SwitchApplicationStateRequest): Promise<void> {
    return apiClient.post<void>(`/projects/${data.project_id}/applications/${data.application_id}/restart`, {});
  },

  async updateGeneral(data: UpdateGeneralSettingsRequest): Promise<void> {
    return apiClient.patch<void>(
      `/projects/${data.project_id}/applications/${data.application_id}/general`,
      data.payload
    );
  },

  async generateDomain(projectId: string, applicationId: string): Promise<{ domain: string }> {
    return apiClient.post<{ domain: string }>(
      `/projects/${projectId}/applications/${applicationId}/domain/generate`,
      {}
    );
  },

  async assignDomain(projectId: string, applicationId: string, domain: string): Promise<void> {
    return apiClient.put<void>(`/projects/${projectId}/applications/${applicationId}/domain`, {
      domain
    });
  },

  async updatePorts(
    projectId: string,
    applicationId: string,
    data: {
      exposed_ports: number[];
      port_mappings: Array<{ host_port: number; container_port: number }>;
    }
  ): Promise<void> {
    return apiClient.put<void>(
      `/projects/${projectId}/applications/${applicationId}/ports`,
      data
    );
  },

  async streamLogs(
    projectId: string,
    applicationId: string,
    follow: boolean = true,
    onLog: (log: string) => void,
    onError?: (error: Error) => void
  ): Promise<() => void> {
    const url = `${apiClient.baseURL}/projects/${projectId}/applications/${applicationId}/logs?follow=${follow}`;
    const abortController = new AbortController();

    fetch(url, {
      method: 'GET',
      headers: apiClient.getHeaders(),
      signal: abortController.signal
    })
      .then(async (response) => {
        if (!response.ok) {
          throw new Error(`Failed to stream logs: ${response.statusText}`);
        }

        const reader = response.body?.getReader();
        const decoder = new TextDecoder();

        if (!reader) {
          throw new Error('Response body is not readable');
        }

        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          const chunk = decoder.decode(value, { stream: true });
          const lines = chunk.split('\n');

          for (const line of lines) {
            if (line.trim()) {
              onLog(line);
            }
          }
        }
      })
      .catch((error) => {
        if (error.name !== 'AbortError') {
          onError?.(error);
        }
      });

    return () => abortController.abort();
  }
};
