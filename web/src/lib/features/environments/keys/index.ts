export const environmentsKeys = {
	all: ['environments'] as const,
	lists: () => [...environmentsKeys.all, 'list'] as const,
	list: (projectId?: string, filters?: Record<string, unknown>) =>
		[...environmentsKeys.lists(), projectId, filters] as const,
	details: () => [...environmentsKeys.all, 'detail'] as const,
	detail: (projectId: string, id: string) =>
		[...environmentsKeys.details(), projectId, id] as const
};
