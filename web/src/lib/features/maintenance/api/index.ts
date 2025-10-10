import { apiClient } from '$lib/api/client';
import type {
	HealthCheckResponse,
	SystemStatusResponse,
	SystemInfoResponse,
	ResourcesResponse,
	DomainListResponse
} from '../types';

export const maintenanceApi = {
	async healthCheck(): Promise<HealthCheckResponse> {
		return apiClient.get<HealthCheckResponse>('/health');
	},

	async systemStatus(): Promise<SystemStatusResponse> {
		return apiClient.get<SystemStatusResponse>('/maintenance/status');
	},

	async systemInfo(): Promise<SystemInfoResponse> {
		return apiClient.get<SystemInfoResponse>('/maintenance/info');
	},

	async getResources(): Promise<ResourcesResponse> {
		return apiClient.get<ResourcesResponse>('/maintenance/resources');
	},

	async listDomains(): Promise<DomainListResponse> {
		return apiClient.get<DomainListResponse>('/maintenance/domains');
	}
};
