import { apiClient } from '$lib/api/client';
import type { Server, CreateServerRequest, UpdateServerRequest } from '../types';

export const serversApi = {
	async list(): Promise<Server[]> {
		return apiClient.get<Server[]>('/servers');
	},

	async get(id: string): Promise<Server> {
		return apiClient.get<Server>(`/servers/${id}`);
	},

	async create(data: CreateServerRequest): Promise<Server> {
		return apiClient.post<Server>('/servers', data);
	},

	async update(id: string, data: UpdateServerRequest): Promise<Server> {
		return apiClient.put<Server>(`/servers/${id}`, data);
	},

	async delete(id: string): Promise<void> {
		return apiClient.delete<void>(`/servers/${id}`);
	},
};
