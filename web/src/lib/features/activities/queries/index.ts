import { createQuery } from '@tanstack/svelte-query';
import { activitiesApi } from '../api';
import { activitiesKeys } from '../keys';

export const createActivitiesQuery = (limit = 20, offset = 0) =>
  createQuery(() => ({
    queryKey: activitiesKeys.recent(limit, offset),
    queryFn: () => activitiesApi.getRecent({ organization_id: "", limit: limit, offset: offset })
  }));

export const createResourceActivitiesQuery = (
  resourceType: string,
  resourceId: string,
  limit = 20,
  offset = 0
) =>
  createQuery(() => ({
    queryKey: activitiesKeys.forResource(resourceType, resourceId, limit, offset),
    queryFn: () => activitiesApi.getForResource
  }));

