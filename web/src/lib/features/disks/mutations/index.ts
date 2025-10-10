import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { disksApi } from '../api';
import { disksKeys } from '../keys';
import type { CreateDiskRequest, ResizeDiskRequest, AttachDiskRequest } from '../types';

export const createDiskMutationQuery = (projectId: string) => {
	const queryClient = useQueryClient();
	return createMutation(() => ({
		mutationFn: (data: CreateDiskRequest) => disksApi.create(projectId, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: disksKeys.lists() });
		}
	}));
};

export const resizeDiskMutationQuery = (projectId: string, diskId: string) => {
	const queryClient = useQueryClient();
	return createMutation(() => ({
		mutationFn: (data: ResizeDiskRequest) => disksApi.resize(projectId, diskId, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: disksKeys.detail(projectId, diskId) });
			queryClient.invalidateQueries({ queryKey: disksKeys.lists() });
		}
	}));
};

export const deleteDiskMutationQuery = (projectId: string, diskId: string) => {
	const queryClient = useQueryClient();
	return createMutation(() => ({
		mutationFn: () => disksApi.delete(projectId, diskId),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: disksKeys.lists() });
			queryClient.invalidateQueries({ queryKey: disksKeys.detail(projectId, diskId) });
		}
	}));
};

export const attachDiskMutationQuery = (projectId: string, diskId: string) => {
	const queryClient = useQueryClient();
	return createMutation(() => ({
		mutationFn: (data: AttachDiskRequest) => disksApi.attach(projectId, diskId, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: disksKeys.detail(projectId, diskId) });
			queryClient.invalidateQueries({ queryKey: disksKeys.lists() });
		}
	}));
};

export const detachDiskMutationQuery = (projectId: string, diskId: string) => {
	const queryClient = useQueryClient();
	return createMutation(() => ({
		mutationFn: () => disksApi.detach(projectId, diskId),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: disksKeys.detail(projectId, diskId) });
			queryClient.invalidateQueries({ queryKey: disksKeys.lists() });
		}
	}));
};
