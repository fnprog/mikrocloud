import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { gitApi } from '../api';
import { gitKeys } from '../keys';
import type {
  ValidateRepositoryRequest,
  ListBranchesRequest,
  DetectBuildMethodRequest,
  CreateGitSourceRequest,
  UpdateGitSourceRequest
} from '../types';

export const validateRepositoryMutationQuery = () =>
  createMutation(() => ({
    mutationFn: (data: ValidateRepositoryRequest) => gitApi.validateRepository(data)
  }));

export const listBranchesMutationQuery = () =>
  createMutation(() => ({
    mutationFn: (data: ListBranchesRequest) => gitApi.listBranches(data)
  }));

export const detectBuildMethodMutationQuery = () =>
  createMutation(() => ({
    mutationFn: (data: DetectBuildMethodRequest) => gitApi.detectBuildMethod(data)
  }));


type CreateGitSourceMutationOptions = {
  onSuccess?: () => void;
  onError?: (error: Error) => void;
};


export const createGitSourceMutationQuery = (options: CreateGitSourceMutationOptions = {}) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: (data: CreateGitSourceRequest) => gitApi.createGitSource(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: gitKeys.sources() });
      options.onSuccess?.()
    },
    onError: (error: Error) => {
      console.error(error)
      options.onError?.(error)
    }
  }));
};

export const updateGitSourceMutationQuery = (sourceId: string, options: CreateGitSourceMutationOptions = {}) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: (data: UpdateGitSourceRequest) => gitApi.updateGitSource(sourceId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: gitKeys.sourceDetail(sourceId) });
      queryClient.invalidateQueries({ queryKey: gitKeys.sources() });
      options.onSuccess?.()
    },
    onError: (error: Error) => {
      console.error(error)
      options.onError?.(error)
    }
  }));
};

export const deleteGitSourceMutationQuery = (options: CreateGitSourceMutationOptions = {}) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: gitApi.deleteGitSource,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: gitKeys.sources() });
      // queryClient.invalidateQueries({ queryKey: gitKeys.sourceDetail(sourceId) });
      options.onSuccess?.()
    },
    onError: (error: Error) => {
      console.error(error)
      options.onError?.(error)
    }
  }));
};
