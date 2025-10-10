import { createQuery } from '@tanstack/svelte-query';
import { environmentsApi } from '../api';
import { environmentsKeys } from '../keys';
import type { Environment } from '../types';

export const createEnvironmentsListQuery = (projectId: string) => {
	return createQuery<Environment[], Error>(() => ({
		queryKey: environmentsKeys.list(projectId),
		queryFn: () => environmentsApi.list(projectId),
		enabled: !!projectId
	}));
};

export const createEnvironmentQuery = (projectId: string, environmentId: string) => {
	return createQuery<Environment, Error>(() => ({
		queryKey: environmentsKeys.detail(projectId, environmentId),
		queryFn: () => environmentsApi.get(projectId, environmentId),
		enabled: !!projectId && !!environmentId
	}));
};
