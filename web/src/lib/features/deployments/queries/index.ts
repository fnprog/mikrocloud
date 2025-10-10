import { createQuery } from '@tanstack/svelte-query';
import { deploymentsApi } from '../api';
import { deploymentsKeys } from '../keys';

export const createDeploymentsListQuery = (projectId: string, applicationId: string) =>
  createQuery(() => ({
    queryKey: deploymentsKeys.list(projectId, applicationId),
    queryFn: () => deploymentsApi.list(projectId, applicationId),
    enabled: !!projectId && !!applicationId,
    refetchInterval: 5000
  }));

export const createDeploymentQuery = (projectId: string, applicationId: string, deploymentId: string) =>
  createQuery(() => ({
    queryKey: deploymentsKeys.detail(projectId, applicationId, deploymentId),
    queryFn: () => deploymentsApi.get(projectId, applicationId, deploymentId)
  }));
