import { apiClient } from '$lib/api/client';
import type { Tunnel, TunnelsResponse, CreateTunnelRequest } from '../types';

export const tunnelsApi = {
	async list(organizationId?: string, projectId?: string): Promise<Tunnel[]> {
		let url = '/tunnels';
		const params = new URLSearchParams();

		if (projectId) {
			params.append('project_id', projectId);
		}

		if (params.toString()) {
			url += `?${params.toString()}`;
		}

		const response = await apiClient.get<TunnelsResponse>(url);
		return response.tunnels;
	},

	async get(tunnelId: string): Promise<Tunnel> {
		return apiClient.get<Tunnel>(`/tunnels/${tunnelId}`);
	},

	async create(data: CreateTunnelRequest): Promise<Tunnel> {
		return apiClient.post<Tunnel>('/tunnels', data);
	},

	async delete(tunnelId: string): Promise<void> {
		return apiClient.delete<void>(`/tunnels/${tunnelId}`);
	},

	async start(tunnelId: string): Promise<Tunnel> {
		return apiClient.post<Tunnel>(`/tunnels/${tunnelId}/start`);
	},

	async stop(tunnelId: string): Promise<Tunnel> {
		return apiClient.post<Tunnel>(`/tunnels/${tunnelId}/stop`);
	},

	async restart(tunnelId: string): Promise<Tunnel> {
		return apiClient.post<Tunnel>(`/tunnels/${tunnelId}/restart`);
	}
};
