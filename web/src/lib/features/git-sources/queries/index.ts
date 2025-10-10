import { createQuery } from '@tanstack/svelte-query';
import { gitApi } from '../api';
import { gitKeys } from '../keys';

export const createGitSourcesListQuery = () =>
	createQuery(() => ({
		queryKey: gitKeys.sourcesList(),
		queryFn: () => gitApi.listGitSources()
	}));

export const createGitSourceQuery = (sourceId: string) =>
	createQuery(() => ({
		queryKey: gitKeys.sourceDetail(sourceId),
		queryFn: () => gitApi.getGitSource(sourceId)
	}));
