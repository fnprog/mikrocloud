import { createQuery } from '@tanstack/svelte-query';
import { proxyApi } from '../api';
import { proxyKeys } from '../keys';

export const createProxyConfigsListQuery = (projectId: string) =>
	createQuery(() => ({
		queryKey: proxyKeys.list(projectId),
		queryFn: () => proxyApi.list(projectId),
		enabled: !!projectId
	}));

export const createProxyConfigDetailQuery = (projectId: string, configId: string) =>
	createQuery(() => ({
		queryKey: proxyKeys.detail(configId),
		queryFn: () => proxyApi.get(projectId, configId),
		enabled: !!projectId && !!configId
	}));
