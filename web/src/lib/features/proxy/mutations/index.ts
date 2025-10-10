import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { proxyApi } from '../api';
import { proxyKeys } from '../keys';
import type { CreateProxyConfigRequest, UpdateProxyConfigRequest } from '../types';

export const createProxyConfigCreateMutation = (projectId: string) => {
	const queryClient = useQueryClient();

	return createMutation(() => ({
		mutationFn: (data: CreateProxyConfigRequest) => proxyApi.create(projectId, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: proxyKeys.list(projectId) });
		}
	}));
};

export const createProxyConfigUpdateMutation = (projectId: string, configId: string) => {
	const queryClient = useQueryClient();

	return createMutation(() => ({
		mutationFn: (data: UpdateProxyConfigRequest) => proxyApi.update(projectId, configId, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: proxyKeys.list(projectId) });
			queryClient.invalidateQueries({ queryKey: proxyKeys.detail(configId) });
		}
	}));
};

export const createProxyConfigDeleteMutation = (projectId: string) => {
	const queryClient = useQueryClient();

	return createMutation(() => ({
		mutationFn: (configId: string) => proxyApi.delete(projectId, configId),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: proxyKeys.list(projectId) });
		}
	}));
};
