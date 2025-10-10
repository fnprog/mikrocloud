import { createQuery } from '@tanstack/svelte-query';
import { organizationsApi } from '../api';
import { organizationsKeys } from '../keys';

export const createOrganizationsListQuery = () =>
  createQuery(() => ({
    queryKey: organizationsKeys.all,
    queryFn: () => organizationsApi.list(),
    refetchOnMount: true,
    staleTime: 5 * 60 * 1000,
  }));

export const createOrganizationDetailQuery = (orgId: string) =>
  createQuery(() => ({
    queryKey: organizationsKeys.detail(orgId),
    queryFn: () => organizationsApi.get(orgId),
    enabled: !!orgId
  }));
