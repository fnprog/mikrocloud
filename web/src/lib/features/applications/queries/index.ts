import { createQuery } from '@tanstack/svelte-query';
import { applicationsApi } from '../api';
import { applicationsKeys } from '../keys';

export const createApplicationsFetchQuery = (projectId: string, environmentId: string) =>
  createQuery(() => ({
    queryKey: applicationsKeys.list(projectId, environmentId),
    queryFn: () => applicationsApi.list(projectId)
  }));

export const createApplicationFetchQuery = (
  projectId: string,
  environmentId: string,
  applicationId: string
) =>
  createQuery(() => ({
    queryKey: applicationsKeys.detail(projectId, environmentId, applicationId),
    queryFn: () => applicationsApi.get(projectId, applicationId),
    enabled: !!projectId && !!applicationId
  }));
