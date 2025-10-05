import { apiClient } from './client';

export type DeploymentStatus = 'pending' | 'building' | 'deploying' | 'success' | 'failed' | 'cancelled';

export interface Deployment {
	id: string;
	application_id: string;
	status: DeploymentStatus;
	commit_hash: string;
	commit_message: string;
	branch: string;
	author: string;
	started_at: string;
	completed_at?: string;
	duration?: number;
	logs?: string;
}

export interface DeploymentsResponse {
	deployments: Deployment[];
}

export const deploymentsApi = {
	async list(projectId: string, applicationId: string): Promise<Deployment[]> {
		const response = await apiClient.get<DeploymentsResponse>(
			`/projects/${projectId}/applications/${applicationId}/deployments`
		);
		return response.deployments;
	},

	async get(projectId: string, applicationId: string, deploymentId: string): Promise<Deployment> {
		return apiClient.get<Deployment>(
			`/projects/${projectId}/applications/${applicationId}/deployments/${deploymentId}`
		);
	},

	async redeploy(projectId: string, applicationId: string, deploymentId: string): Promise<Deployment> {
		return apiClient.post<Deployment>(
			`/projects/${projectId}/applications/${applicationId}/deployments/${deploymentId}/redeploy`,
			{}
		);
	},

	async cancel(projectId: string, applicationId: string, deploymentId: string): Promise<void> {
		return apiClient.post<void>(
			`/projects/${projectId}/applications/${applicationId}/deployments/${deploymentId}/cancel`,
			{}
		);
	}
};
