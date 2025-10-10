import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { databasesApi } from '../api';
import { databasesKeys } from '../keys';
import type { CreateDatabaseRequest } from '../types';


type GeneralMutationOptions = {
  onSuccess?: () => void;
  onError?: (error: Error) => void;
};


export const createDatabaseMutationQuery = (projectId: string, options: GeneralMutationOptions = {}) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: (data: CreateDatabaseRequest) => databasesApi.create(projectId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: databasesKeys.lists() });
      options.onSuccess?.();
    },
    onError: (error: Error) => {
      console.error(error);
      options.onError?.(error);
    }
  }));
};

export const deleteDatabaseMutationQuery = (projectId: string, databaseId: string) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: () => databasesApi.delete(projectId, databaseId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: databasesKeys.lists() });
      queryClient.invalidateQueries({ queryKey: databasesKeys.detail(projectId, databaseId) });
    }
  }));
};

export const startDatabaseMutationQuery = (projectId: string, databaseId: string) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: () => databasesApi.start(projectId, databaseId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: databasesKeys.detail(projectId, databaseId) });
      queryClient.invalidateQueries({ queryKey: databasesKeys.lists() });
    }
  }));
};

export const stopDatabaseMutationQuery = (projectId: string, databaseId: string) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: () => databasesApi.stop(projectId, databaseId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: databasesKeys.detail(projectId, databaseId) });
      queryClient.invalidateQueries({ queryKey: databasesKeys.lists() });
    }
  }));
};

export const restartDatabaseMutationQuery = (projectId: string, databaseId: string) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: () => databasesApi.restart(projectId, databaseId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: databasesKeys.detail(projectId, databaseId) });
      queryClient.invalidateQueries({ queryKey: databasesKeys.lists() });
    }
  }));
};
