import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { environmentsApi } from '../api';
import { environmentsKeys } from '../keys';
import type { Environment, CreateEnvironmentRequest } from '../types';

type CreateEnvironmentMutationOptions = {
	onSuccess?: () => void;
	onError?: (error: Error) => void;
};

export const createEnvironmentMutationQuery = (
	options: CreateEnvironmentMutationOptions = {}
) => {
	const queryClient = useQueryClient();

	return createMutation<
		Environment,
		Error,
		{ projectId: string; data: CreateEnvironmentRequest }
	>(() => ({
		mutationFn: ({ projectId, data }) => environmentsApi.create(projectId, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: environmentsKeys.all });
			options.onSuccess?.();
		},
		onError: (error: Error) => {
			console.error(error);
			options.onError?.(error);
		}
	}));
};

type UpdateEnvironmentMutationOptions = {
	onSuccess?: () => void;
	onError?: (error: Error) => void;
};

export const updateEnvironmentMutationQuery = (
	options: UpdateEnvironmentMutationOptions = {}
) => {
	const queryClient = useQueryClient();

	return createMutation<
		Environment,
		Error,
		{ projectId: string; environmentId: string; data: Partial<CreateEnvironmentRequest> }
	>(() => ({
		mutationFn: ({ projectId, environmentId, data }) =>
			environmentsApi.update(projectId, environmentId, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: environmentsKeys.all });
			options.onSuccess?.();
		},
		onError: (error: Error) => {
			console.error(error);
			options.onError?.(error);
		}
	}));
};

type DeleteEnvironmentMutationOptions = {
	onSuccess?: () => void;
	onError?: (error: Error) => void;
};

export const deleteEnvironmentMutationQuery = (
	options: DeleteEnvironmentMutationOptions = {}
) => {
	const queryClient = useQueryClient();

	return createMutation<void, Error, { projectId: string; environmentId: string }>(() => ({
		mutationFn: ({ projectId, environmentId }) =>
			environmentsApi.delete(projectId, environmentId),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: environmentsKeys.all });
			options.onSuccess?.();
		},
		onError: (error: Error) => {
			console.error(error);
			options.onError?.(error);
		}
	}));
};
