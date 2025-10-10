import { apiClient } from '$lib/api/client';
import type { Organization } from '../types';

export const organizationsApi = {
	async list(): Promise<Organization[]> {
		return apiClient.get<Organization[]>('/organizations');
	},

	async get(orgId: string): Promise<Organization> {
		return apiClient.get<Organization>(`/organizations/${orgId}`);
	}
};
