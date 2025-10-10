export const disksKeys = {
	all: ['disks'] as const,
	lists: () => [...disksKeys.all, 'list'] as const,
	list: (projectId: string) => [...disksKeys.lists(), { projectId }] as const,
	details: () => [...disksKeys.all, 'detail'] as const,
	detail: (projectId: string, diskId: string) =>
		[...disksKeys.details(), { projectId, diskId }] as const
};
