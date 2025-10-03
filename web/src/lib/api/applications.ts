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
	}
};
