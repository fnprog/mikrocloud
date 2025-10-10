export const gitKeys = {
	all: ['git'] as const,
	sources: () => [...gitKeys.all, 'sources'] as const,
	sourcesList: () => [...gitKeys.sources(), 'list'] as const,
	sourceDetail: (sourceId: string) => [...gitKeys.sources(), 'detail', { sourceId }] as const
};
