import { createQuery } from '@tanstack/svelte-query';
import { projectsApi } from '../api';
import { projectsKeys } from '../keys';
import type { Project } from '../types';

export const createProjectsQuery = () =>
  createQuery(() => ({
    queryKey: projectsKeys.all,
    queryFn: () => projectsApi.list()
  }));

export const createProjectsListQuery = () => {
  return createQuery(() => ({
    queryKey: projectsKeys.lists(),
    queryFn: projectsApi.list
  }));
};

export const createProjectQuery = (projectId: string) => {
  return createQuery<Project, Error>(() => ({
    queryKey: projectsKeys.detail(projectId),
    queryFn: () => projectsApi.get(projectId),
    enabled: !!projectId
  }));
};
