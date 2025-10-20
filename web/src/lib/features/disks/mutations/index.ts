import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { disksApi } from '../api';
import { disksKeys } from '../keys';
import type { CreateDiskRequest, ResizeDiskRequest, AttachDiskRequest } from '../types';

export const createDiskMutationQuery = (projectId: string) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: (data: CreateDiskRequest) => disksApi.create(projectId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: disksKeys.lists() });
    }
  }));
};

export const resizeDiskMutationQuery = (projectId: string) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: (data: ResizeDiskRequest) => disksApi.resize(projectId, data.disk_id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: disksKeys.list(projectId) });
      queryClient.invalidateQueries({ queryKey: disksKeys.lists() });
    }
  }));
};

export const deleteDiskMutationQuery = (projectId: string) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: (diskId: string) => disksApi.delete(projectId, diskId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: disksKeys.lists() });
      queryClient.invalidateQueries({ queryKey: disksKeys.list(projectId) });
    }
  }));
};

export const attachDiskMutationQuery = (projectId: string) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: (data: AttachDiskRequest) => disksApi.attach(projectId, data.disk_id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: disksKeys.list(projectId) });
      queryClient.invalidateQueries({ queryKey: disksKeys.lists() });
    }
  }));
};

export const detachDiskMutationQuery = (projectId: string) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: (diskId: string) => disksApi.detach(projectId, diskId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: disksKeys.list(projectId) });
      queryClient.invalidateQueries({ queryKey: disksKeys.lists() });
    }
  }));
};
