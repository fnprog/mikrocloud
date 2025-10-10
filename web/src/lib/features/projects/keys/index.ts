export const projectsKeys = {
	all: ['projects'] as const,
	lists: () => [...projectsKeys.all, 'list'] as const,
	list: (filters?: Record<string, unknown>) => [...projectsKeys.lists(), filters] as const,
	details: () => [...projectsKeys.all, 'detail'] as const,
	detail: (id: string) => [...projectsKeys.details(), id] as const
};
