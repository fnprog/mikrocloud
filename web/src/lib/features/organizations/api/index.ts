import { apiClient } from '$lib/api/client';
import type { Organization } from '../types';

export interface CreateOrganizationRequest {
	name: string;
	slug: string;
	description?: string;
	billing_email?: string;
}

export interface UpdateOrganizationRequest {
	name?: string;
	description?: string;
	billing_email?: string;
}

export const organizationsApi = {
	async list(): Promise<Organization[]> {
		return apiClient.get<Organization[]>('/organizations');
	},

	async get(orgId: string): Promise<Organization> {
		return apiClient.get<Organization>(`/organizations/${orgId}`);
	},

	async create(data: CreateOrganizationRequest): Promise<Organization> {
		return apiClient.post<Organization>('/organizations', data);
	},

	async update(orgId: string, data: UpdateOrganizationRequest): Promise<Organization> {
		return apiClient.put<Organization>(`/organizations/${orgId}`, data);
	},

	async delete(orgId: string): Promise<void> {
		return apiClient.delete(`/organizations/${orgId}`);
	}
};
