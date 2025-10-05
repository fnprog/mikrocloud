import { apiClient } from './client';

export type ApplicationStatus = 'pending' | 'running' | 'stopped' | 'failed' | 'deploying';

export interface Application {
	id: string;
	name: string;
	description: string;
	project_id: string;
	environment_id: string;
	domain: string;
	status: ApplicationStatus;
	created_at: string;
}

export interface CreateApplicationRequest {
	name: string;
	description?: string;
	environment_id: string;
	deployment_source: {
		type: 'git' | 'docker';
		config: any;
	};
	buildpack: {
		type: string;
		config: any;
	};
	env_vars?: Record<string, string>;
}

export interface ApplicationsResponse {
	applications: Application[];
}

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

	async create(projectId: string, data: CreateApplicationRequest): Promise<Application> {
		return apiClient.post<Application>(`/projects/${projectId}/applications`, data);
	},

	async delete(projectId: string, applicationId: string): Promise<void> {
		return apiClient.delete<void>(`/projects/${projectId}/applications/${applicationId}`);
	},

	async start(projectId: string, applicationId: string): Promise<void> {
		return apiClient.post<void>(`/projects/${projectId}/applications/${applicationId}/start`, {});
	},

	async stop(projectId: string, applicationId: string): Promise<void> {
		return apiClient.post<void>(`/projects/${projectId}/applications/${applicationId}/stop`, {});
	},

	async restart(projectId: string, applicationId: string): Promise<void> {
		return apiClient.post<void>(`/projects/${projectId}/applications/${applicationId}/restart`, {});
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
