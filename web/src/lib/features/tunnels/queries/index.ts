import { createQuery } from '@tanstack/svelte-query';
import { tunnelsApi } from '../api';
import { tunnelsKeys } from '../keys';

export const createTunnelsFetchQuery = (organizationId?: string, projectId?: string) =>
	createQuery(() => ({
		queryKey: tunnelsKeys.list(organizationId, projectId),
		queryFn: () => tunnelsApi.list(organizationId, projectId)
	}));

export const createTunnelFetchQuery = (tunnelId: string) =>
	createQuery(() => ({
		queryKey: tunnelsKeys.detail(tunnelId),
		queryFn: () => tunnelsApi.get(tunnelId),
		enabled: !!tunnelId
	}));
