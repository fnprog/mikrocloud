import type { DatabaseType } from '../types';

export const databasesKeys = {
	all: ['databases'] as const,
	lists: () => [...databasesKeys.all, 'list'] as const,
	list: (projectId: string, environmentId?: string) =>
		[...databasesKeys.lists(), { projectId, environmentId }] as const,
	details: () => [...databasesKeys.all, 'detail'] as const,
	detail: (projectId: string, databaseId: string) =>
		[...databasesKeys.details(), { projectId, databaseId }] as const,
	types: (projectId: string) => [...databasesKeys.all, 'types', { projectId }] as const,
	defaultConfig: (projectId: string, type: DatabaseType) =>
		[...databasesKeys.all, 'defaultConfig', { projectId, type }] as const
};
