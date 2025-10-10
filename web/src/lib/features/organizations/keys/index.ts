export const organizationsKeys = {
	all: ['organizations'] as const,
	lists: () => [...organizationsKeys.all, 'list'] as const,
	list: () => [...organizationsKeys.lists()] as const,
	details: () => [...organizationsKeys.all, 'detail'] as const,
	detail: (id: string) => [...organizationsKeys.details(), id] as const
};
