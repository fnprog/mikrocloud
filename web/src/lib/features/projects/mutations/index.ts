import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { projectsApi } from '../api';
import { projectsKeys } from '../keys';
import type { Project, CreateProjectRequest } from '../types';

type CreateProjectMutationOptions = {
	onSuccess?: () => void;
	onError?: (error: Error) => void;
};

export const createProjectMutationQuery = (options: CreateProjectMutationOptions = {}) => {
	const queryClient = useQueryClient();

	return createMutation<Project, Error, CreateProjectRequest>(() => {
		return {
			mutationFn: projectsApi.create,
			onSuccess: () => {
				queryClient.invalidateQueries({ queryKey: projectsKeys.all });
				options.onSuccess?.();
			},
			onError: (error: Error) => {
				console.error(error);
				options.onError?.(error);
			}
		};
	});
};

type UpdateProjectMutationOptions = {
	onSuccess?: () => void;
	onError?: (error: Error) => void;
};

export const updateProjectMutationQuery = (options: UpdateProjectMutationOptions = {}) => {
	const queryClient = useQueryClient();

	return createMutation<
		Project,
		Error,
		{ id: string; data: Partial<CreateProjectRequest> }
	>(() => ({
		mutationFn: ({ id, data }) => projectsApi.update(id, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: projectsKeys.all });
			options.onSuccess?.();
		},
		onError: (error: Error) => {
			console.error(error);
			options.onError?.(error);
		}
	}));
};

type deleteProjectMutationOptions = {
	onSuccess?: () => void;
	onError?: (error: Error) => void;
};

export const deleteProjectMutationQuery = (options: deleteProjectMutationOptions = {}) => {
	const queryClient = useQueryClient();

	return createMutation<void, Error, string>(() => ({
		mutationFn: projectsApi.delete,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: projectsKeys.all });
			options.onSuccess?.();
		},
		onError: (error: Error) => {
			console.error(error);
			options.onError?.(error);
		}
	}));
};
