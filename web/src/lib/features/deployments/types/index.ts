export type DeploymentStatus = 'pending' | 'queued' | 'building' | 'deploying' | 'running' | 'success' | 'failed' | 'cancelled' | 'stopped';

export type LogLevel = 'info' | 'warn' | 'error' | 'success';

export interface StructuredLog {
  timestamp: string;
  relative_seconds: number;
  level: LogLevel;
  message: string;
  raw: string;
}

export interface Deployment {
  id: string;
  application_id: string;
  deployment_number: number;
  status: DeploymentStatus;
  is_production: boolean;
  commit_hash: string;
  commit_message: string;
  branch: string;
  author: string;
  build_started_at: string;
  build_completed_at?: string;
  build_duration_seconds?: number;
  triggered_by_username?: string;
  build_logs?: string;
  error_message?: string;
  created_at: string;
  updated_at: string;
}

export interface DeploymentsResponse {
  deployments: Deployment[];
}
