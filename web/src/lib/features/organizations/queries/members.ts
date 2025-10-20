import { createQuery } from '@tanstack/svelte-query';
import * as membersApi from '$lib/api/organizations/members';
import { membersKeys } from '../keys/members';

export const createMembersListQuery = (organizationId: string) =>
	createQuery(() => ({
		queryKey: membersKeys.list(organizationId),
		queryFn: () => membersApi.listMembers(organizationId),
		enabled: !!organizationId,
		refetchOnMount: true,
		staleTime: 2 * 60 * 1000
	}));
