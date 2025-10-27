import { createQuery } from '@tanstack/svelte-query';
import { authApi } from '../api';
import { authKeys } from '../keys';

export const createProfileQuery = () =>
	createQuery(() => ({
		queryKey: authKeys.profile(),
		queryFn: () => authApi.getProfile(),
		enabled: authApi.isAuthenticated(),
		staleTime: 1000 * 60 * 5
	}));

export const createSetupStatusQuery = () =>
	createQuery(() => ({
		queryKey: authKeys.setupStatus(),
		queryFn: () => authApi.getSetupStatus(),
		staleTime: 1000 * 60 * 5
	}));
