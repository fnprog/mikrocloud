export const applicationsKeys = {
  all: ['applications'] as const,

  lists: () => [...applicationsKeys.all, 'list'] as const,
  list: (projectId: string, environmentId: string) =>
    [...applicationsKeys.lists(), { projectId, environmentId }] as const,
  details: () => [...applicationsKeys.all, 'detail'] as const,
  detail: (projectId: string, environmentId: string, applicationId: string) =>
    [...applicationsKeys.details(), projectId, environmentId, applicationId] as const
};
