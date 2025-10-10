import { createQuery } from '@tanstack/svelte-query';
import { maintenanceApi } from '../api';
import { maintenanceKeys } from '../keys';

export const createHealthCheckQuery = () =>
	createQuery(() => ({
		queryKey: maintenanceKeys.health(),
		queryFn: () => maintenanceApi.healthCheck(),
		staleTime: 1000 * 30,
		refetchInterval: 1000 * 60
	}));

export const createSystemStatusQuery = () =>
	createQuery(() => ({
		queryKey: maintenanceKeys.systemStatus(),
		queryFn: () => maintenanceApi.systemStatus(),
		staleTime: 1000 * 30,
		refetchInterval: 1000 * 60
	}));

export const createSystemInfoQuery = () =>
	createQuery(() => ({
		queryKey: maintenanceKeys.systemInfo(),
		queryFn: () => maintenanceApi.systemInfo(),
		staleTime: 1000 * 60 * 5
	}));

export const createResourcesQuery = () =>
	createQuery(() => ({
		queryKey: maintenanceKeys.resources(),
		queryFn: () => maintenanceApi.getResources(),
		staleTime: 1000 * 60
	}));

export const createDomainsListQuery = () =>
	createQuery(() => ({
		queryKey: maintenanceKeys.domains(),
		queryFn: () => maintenanceApi.listDomains()
	}));
