import { createQuery } from '@tanstack/svelte-query';
import { serversApi } from '../api';
import { serversKeys } from '../keys';

export const createServersListQuery = () =>
	createQuery(() => ({
		queryKey: serversKeys.list(),
		queryFn: () => serversApi.list()
	}));

export const createServerQuery = (serverId: string) =>
	createQuery(() => ({
		queryKey: serversKeys.detail(serverId),
		queryFn: () => serversApi.get(serverId)
	}));
