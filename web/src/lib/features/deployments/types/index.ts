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
