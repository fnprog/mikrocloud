import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { deploymentsApi } from '../api';
import { deploymentsKeys } from '../keys';

export const redeployDeploymentMutationQuery = (projectId: string, applicationId: string, deploymentId: string) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: () => deploymentsApi.redeploy(projectId, applicationId, deploymentId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deploymentsKeys.lists() });
      queryClient.invalidateQueries({ queryKey: deploymentsKeys.detail(projectId, applicationId, deploymentId) });
    }
  }));
};

export const cancelDeploymentMutationQuery = (projectId: string, applicationId: string, deploymentId: string) => {
  const queryClient = useQueryClient();
  return createMutation(() => ({
    mutationFn: () => deploymentsApi.cancel(projectId, applicationId, deploymentId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deploymentsKeys.lists() });
      queryClient.invalidateQueries({ queryKey: deploymentsKeys.detail(projectId, applicationId, deploymentId) });
    }
  }));
};
