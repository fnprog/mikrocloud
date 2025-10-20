import { createQuery } from '@tanstack/svelte-query';
import { disksApi } from '../api';
import { disksKeys } from '../keys';

export const createDisksListQuery = (projectId: string) =>
  createQuery(() => ({
    queryKey: disksKeys.list(projectId),
    queryFn: () => disksApi.list(projectId),
    enabled: !!projectId
  }));

export const createDiskQuery = (projectId: string, diskId: string) =>
  createQuery(() => ({
    queryKey: disksKeys.detail(projectId, diskId),
    queryFn: () => disksApi.get(projectId, diskId)
  }));
