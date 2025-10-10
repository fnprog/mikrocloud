export const deploymentsKeys = {
  all: ['deployments'] as const,
  lists: () => [...deploymentsKeys.all, 'list'] as const,
  list: (projectId: string, applicationId: string) =>
    [...deploymentsKeys.lists(), { projectId, applicationId }] as const,
  details: () => [...deploymentsKeys.all, 'detail'] as const,
  detail: (projectId: string, applicationId: string, deploymentId: string) =>
    [...deploymentsKeys.details(), { projectId, applicationId, deploymentId }] as const
};
