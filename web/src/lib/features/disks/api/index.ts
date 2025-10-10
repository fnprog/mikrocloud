import { apiClient } from '$lib/api/client';
import type { Disk, CreateDiskRequest, ResizeDiskRequest, AttachDiskRequest, ListDisksResponse } from '../types';

export const disksApi = {
	async list(projectId: string): Promise<Disk[]> {
		const response = await apiClient.get<ListDisksResponse>(
			`/projects/${projectId}/disks`
		);
		return response.disks;
	},

	async get(projectId: string, diskId: string): Promise<Disk> {
		return apiClient.get<Disk>(`/projects/${projectId}/disks/${diskId}`);
	},

	async create(projectId: string, data: CreateDiskRequest): Promise<Disk> {
		return apiClient.post<Disk>(`/projects/${projectId}/disks`, data);
	},

	async resize(projectId: string, diskId: string, data: ResizeDiskRequest): Promise<Disk> {
		return apiClient.put<Disk>(`/projects/${projectId}/disks/${diskId}/resize`, data);
	},

	async delete(projectId: string, diskId: string): Promise<void> {
		return apiClient.delete<void>(`/projects/${projectId}/disks/${diskId}`);
	},

	async attach(projectId: string, diskId: string, data: AttachDiskRequest): Promise<Disk> {
		return apiClient.post<Disk>(`/projects/${projectId}/disks/${diskId}/attach`, data);
	},

	async detach(projectId: string, diskId: string): Promise<Disk> {
		return apiClient.post<Disk>(`/projects/${projectId}/disks/${diskId}/detach`, {});
	},
};
