import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { serversApi } from '../api';
import { serversKeys } from '../keys';
import type { CreateServerRequest, UpdateServerRequest } from '../types';

export const createServerMutationQuery = () => {
	const queryClient = useQueryClient();
	return createMutation(() => ({
		mutationFn: (data: CreateServerRequest) => serversApi.create(data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: serversKeys.lists() });
		}
	}));
};

export const updateServerMutationQuery = (serverId: string) => {
	const queryClient = useQueryClient();
	return createMutation(() => ({
		mutationFn: (data: UpdateServerRequest) => serversApi.update(serverId, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: serversKeys.detail(serverId) });
			queryClient.invalidateQueries({ queryKey: serversKeys.lists() });
		}
	}));
};

export const deleteServerMutationQuery = (serverId: string) => {
	const queryClient = useQueryClient();
	return createMutation(() => ({
		mutationFn: () => serversApi.delete(serverId),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: serversKeys.lists() });
			queryClient.invalidateQueries({ queryKey: serversKeys.detail(serverId) });
		}
	}));
};
