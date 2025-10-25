import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { tunnelsApi } from '../api';
import { tunnelsKeys } from '../keys';
import type { CreateTunnelRequest } from '../types';

type GeneralMutationOptions = {
  onSuccess?: () => void;
  onError?: (error: Error) => void;
};

export const createTunnelMutationQuery = (options: GeneralMutationOptions = {}) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: (data: CreateTunnelRequest) => tunnelsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tunnelsKeys.lists() });
      options.onSuccess?.();
    },
    onError: (error: Error) => {
      console.error(error);
      options.onError?.(error);
    }
  }));
};

export const deleteTunnelMutationQuery = (tunnelId: string, options: GeneralMutationOptions = {}) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: () => tunnelsApi.delete(tunnelId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tunnelsKeys.lists() });
      queryClient.invalidateQueries({ queryKey: tunnelsKeys.detail(tunnelId) });
      options.onSuccess?.();
    },
    onError: (error: Error) => {
      console.error(error);
      options.onError?.(error);
    }
  }));
};

export const startTunnelMutationQuery = (options: GeneralMutationOptions = {}) => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: tunnelsApi.start,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tunnelsKeys.lists() });
      options.onSuccess?.();
    },
    onError: (error: Error) => {
      console.error(error);
      options.onError?.(error);
    }
  }));
};

export const stopTunnelMutationQuery = (tunnelId: string, options: GeneralMutationOptions = {}) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: () => tunnelsApi.stop(tunnelId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tunnelsKeys.detail(tunnelId) });
      queryClient.invalidateQueries({ queryKey: tunnelsKeys.lists() });
      options.onSuccess?.();
    },
    onError: (error: Error) => {
      console.error(error);
      options.onError?.(error);
    }
  }));
};

export const restartTunnelMutationQuery = (
  tunnelId: string,
  options: GeneralMutationOptions = {}
) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: () => tunnelsApi.restart(tunnelId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tunnelsKeys.detail(tunnelId) });
      queryClient.invalidateQueries({ queryKey: tunnelsKeys.lists() });
      options.onSuccess?.();
    },
    onError: (error: Error) => {
      console.error(error);
      options.onError?.(error);
    }
  }));
};
