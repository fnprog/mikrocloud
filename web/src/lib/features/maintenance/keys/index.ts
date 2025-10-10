export const maintenanceKeys = {
	all: ['maintenance'] as const,
	health: () => [...maintenanceKeys.all, 'health'] as const,
	systemStatus: () => [...maintenanceKeys.all, 'system-status'] as const,
	systemInfo: () => [...maintenanceKeys.all, 'system-info'] as const,
	resources: () => [...maintenanceKeys.all, 'resources'] as const,
	domains: () => [...maintenanceKeys.all, 'domains'] as const
};
