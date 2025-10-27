import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { authApi } from '../api';
import { authKeys } from '../keys';
import type { LoginRequest, RegisterRequest, UpdateProfileRequest } from '../types';

type GenericMutationOptions = {
  onSuccess?: () => void;
  onError?: (error: Error) => void;
};



export const createLoginMutation = (options: GenericMutationOptions = {}) => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: (credentials: LoginRequest) => authApi.login(credentials),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: authKeys.profile() });
      options.onSuccess?.()
    },
    onError: (error: Error) => {
      options.onError?.(error)
    }
  }));
};

export const createUpdateAccountMutation = (options: GenericMutationOptions = {}) => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: (data: UpdateProfileRequest) => authApi.updateProfile(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: authKeys.profile() });
      options.onSuccess?.()
    },

    onError: (error: Error) => {
      options.onError?.(error)
    }
  }));
};


export const createDeleteAccountMutation = (options: GenericMutationOptions = {}) => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: authApi.deleteProfile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: authKeys.profile() });
      options.onSuccess?.()
    },

    onError: (error: Error) => {
      options.onError?.(error)
    }
  }));
};



export const createRegisterMutation = (options: GenericMutationOptions = {}) => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: (data: RegisterRequest) => authApi.register(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: authKeys.profile() });
      options.onSuccess?.()
    },

    onError: (error: Error) => {
      options.onError?.(error)
    }
  }));
};

export const createLogoutMutation = (options: GenericMutationOptions = {}) => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: () => {
      authApi.logout();
      return Promise.resolve();
    },
    onSuccess: () => {
      queryClient.clear();
      options.onSuccess?.()
    },
    onError: (error: Error) => {
      options.onError?.(error)
    }
  }));
};

export const createRequestPasswordResetMutation = (options: GenericMutationOptions = {}) => {
  return createMutation(() => ({
    mutationFn: (data: { email: string }) => authApi.requestPasswordReset(data),
    onSuccess: () => {
      options.onSuccess?.()
    },
    onError: (error: Error) => {
      options.onError?.(error)
    }
  }));
};
