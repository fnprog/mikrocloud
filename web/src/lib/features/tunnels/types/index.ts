export type TunnelStatus = 'stopped' | 'starting' | 'running' | 'stopping' | 'error';
export type HealthStatus = 'healthy' | 'unhealthy' | 'unknown';

export interface Tunnel {
	id: string;
	name: string;
	project_id?: string;
	organization_id: string;
	status: TunnelStatus;
	health_status: HealthStatus;
	container_id: string;
	last_error?: string;
	created_at: string;
	updated_at: string;
}

export interface CreateTunnelRequest {
	name: string;
	project_id?: string;
	tunnel_token: string;
}

export interface TunnelsResponse {
	tunnels: Tunnel[];
}
