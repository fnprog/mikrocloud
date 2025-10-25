export const tunnelsKeys = {
	all: ['tunnels'] as const,
	lists: () => [...tunnelsKeys.all, 'list'] as const,
	list: (organizationId?: string, projectId?: string) =>
		[...tunnelsKeys.lists(), { organizationId, projectId }] as const,
	details: () => [...tunnelsKeys.all, 'detail'] as const,
	detail: (tunnelId: string) => [...tunnelsKeys.details(), { tunnelId }] as const
};
