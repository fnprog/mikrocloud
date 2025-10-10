export const serversKeys = {
	all: ['servers'] as const,
	lists: () => [...serversKeys.all, 'list'] as const,
	list: () => [...serversKeys.lists()] as const,
	details: () => [...serversKeys.all, 'detail'] as const,
	detail: (serverId: string) => [...serversKeys.details(), { serverId }] as const
};
