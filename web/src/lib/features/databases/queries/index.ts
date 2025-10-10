import { createQuery } from '@tanstack/svelte-query';
import { databasesApi } from '../api';
import { databasesKeys } from '../keys';
import type { DatabaseType } from '../types';

export const createDatabasesFetchQuery = (projectId: string, environmentId?: string) =>
  createQuery(() => ({
    queryKey: databasesKeys.list(projectId, environmentId),
    queryFn: () => databasesApi.list(projectId, environmentId)
  }));

export const createDatabaseFetchQuery = (projectId: string, databaseId: string) =>
  createQuery(() => ({
    queryKey: databasesKeys.detail(projectId, databaseId),
    queryFn: () => databasesApi.get(projectId, databaseId),
    enabled: !!projectId && !!databaseId
  }));

export const createDatabasesTypeFetchQuery = (projectId: string) =>
  createQuery(() => ({
    queryKey: databasesKeys.types(projectId),
    queryFn: () => databasesApi.getTypes(projectId)
  }));

export const createDatabaseDefaultConfigFetchQuery = (projectId: string, type: DatabaseType) =>
  createQuery(() => ({
    queryKey: databasesKeys.defaultConfig(projectId, type),
    queryFn: () => databasesApi.getDefaultConfig(projectId, type)
  }));


