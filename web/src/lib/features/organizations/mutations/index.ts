import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { organizationsApi, type CreateOrganizationRequest, type UpdateOrganizationRequest } from '../api';
import { organizationsKeys } from '../keys';
import type { Organization } from '../types';

type CreateOrganizationMutationOptions = {
	onSuccess?: () => void;
	onError?: (error: Error) => void;
};

export const createOrganizationMutationQuery = (options: CreateOrganizationMutationOptions = {}) => {
	const queryClient = useQueryClient();

	return createMutation<Organization, Error, CreateOrganizationRequest>(() => {
		return {
			mutationFn: organizationsApi.create,
			onSuccess: () => {
				queryClient.invalidateQueries({ queryKey: organizationsKeys.all });
				options.onSuccess?.();
			},
			onError: (error: Error) => {
				console.error(error);
				options.onError?.(error);
			}
		};
	});
};

type UpdateOrganizationMutationOptions = {
	onSuccess?: () => void;
	onError?: (error: Error) => void;
};

export const updateOrganizationMutationQuery = (options: UpdateOrganizationMutationOptions = {}) => {
	const queryClient = useQueryClient();

	return createMutation<
		Organization,
		Error,
		{ id: string; data: UpdateOrganizationRequest }
	>(() => ({
		mutationFn: ({ id, data }) => organizationsApi.update(id, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: organizationsKeys.all });
			options.onSuccess?.();
		},
		onError: (error: Error) => {
			console.error(error);
			options.onError?.(error);
		}
	}));
};

type DeleteOrganizationMutationOptions = {
	onSuccess?: () => void;
	onError?: (error: Error) => void;
};

export const deleteOrganizationMutationQuery = (options: DeleteOrganizationMutationOptions = {}) => {
	const queryClient = useQueryClient();

	return createMutation<void, Error, string>(() => ({
		mutationFn: organizationsApi.delete,
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: organizationsKeys.all });
			options.onSuccess?.();
		},
		onError: (error: Error) => {
			console.error(error);
			options.onError?.(error);
		}
	}));
};
