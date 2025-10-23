import { apiClient } from '$lib/api/client';
import type { Deployment, DeploymentsResponse, DeploymentStatus, StructuredLog } from '../types';

export type { DeploymentStatus } from '../types';

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
	},

	async getLogs(
		projectId: string,
		applicationId: string,
		deploymentId: string,
		type: 'build' | 'deploy' | 'all' = 'all'
	): Promise<StructuredLog[] | string> {
		const response = await apiClient.get<{ logs: StructuredLog[] | string; type: string }>(
			`/projects/${projectId}/applications/${applicationId}/deployments/${deploymentId}/logs?type=${type}`
		);
		return response.logs;
	},

	streamLogs(
		projectId: string,
		applicationId: string,
		deploymentId: string,
		onLog: (log: StructuredLog) => void,
		onDone: (status: DeploymentStatus) => void,
		onError: (error: Error) => void
	): () => void {
		const token = localStorage.getItem('auth_token');
		const url = `${apiClient.baseURL}/projects/${projectId}/applications/${applicationId}/deployments/${deploymentId}/logs/stream`;

		const eventSource = new EventSource(`${url}?token=${token}`);

		eventSource.onmessage = (event) => {
			try {
				const log = JSON.parse(event.data) as StructuredLog;
				onLog(log);
			} catch (e) {
				console.error('Failed to parse log:', e);
			}
		};

		eventSource.addEventListener('done', (event) => {
			const data = JSON.parse(event.data);
			onDone(data.status);
			eventSource.close();
		});

		eventSource.onerror = (error) => {
			onError(new Error('Failed to connect to log stream'));
			eventSource.close();
		};

		return () => eventSource.close();
	}
};
