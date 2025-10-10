import { apiClient } from '$lib/api/client';
import type { ProxyConfig, CreateProxyConfigRequest, UpdateProxyConfigRequest, ProxyConfigListResponse } from '../types';

export const proxyApi = {
	async list(projectId: string): Promise<ProxyConfigListResponse> {
		return apiClient.get<ProxyConfigListResponse>(`/projects/${projectId}/proxy`);
	},

	async get(projectId: string, configId: string): Promise<ProxyConfig> {
		return apiClient.get<ProxyConfig>(`/projects/${projectId}/proxy/${configId}`);
	},

	async create(projectId: string, data: CreateProxyConfigRequest): Promise<ProxyConfig> {
		return apiClient.post<ProxyConfig>(`/projects/${projectId}/proxy`, data);
	},

	async update(projectId: string, configId: string, data: UpdateProxyConfigRequest): Promise<ProxyConfig> {
		return apiClient.put<ProxyConfig>(`/projects/${projectId}/proxy/${configId}`, data);
	},

	async delete(projectId: string, configId: string): Promise<void> {
		return apiClient.delete(`/projects/${projectId}/proxy/${configId}`);
	}
};
