export const proxyKeys = {
	all: ['proxy'] as const,
	lists: () => [...proxyKeys.all, 'list'] as const,
	list: (projectId?: string) => [...proxyKeys.lists(), projectId] as const,
	details: () => [...proxyKeys.all, 'detail'] as const,
	detail: (id: string) => [...proxyKeys.details(), id] as const
};
