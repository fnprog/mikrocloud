import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { applicationsApi } from '../api';
import { applicationsKeys } from '../keys';

type CreateApplicationMutationOptions = {
  onSuccess?: () => void;
  onError?: (error: Error) => void;
};

//TODO: Invalidating all is a bit too much, we should only invalidate the project env

export const createApplicationMutationQuery = (
  options: CreateApplicationMutationOptions = {}
) => {
  const queryClient = useQueryClient();

  return createMutation(() => {
    return {
      mutationFn: applicationsApi.create,
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: applicationsKeys.all });
        options.onSuccess?.();
      },
      onError: (error: Error) => {
        console.error(error);
        options.onError?.(error);
      }
    };
  });
};

type deleteApplicationMutationOptions = {
  onSuccess?: () => void;
  onError?: (error: Error) => void;
};

export const deleteApplicationMutationQuery = (
  options: deleteApplicationMutationOptions = {}
) => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: applicationsApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: applicationsKeys.all });
      options.onSuccess?.();
    },
    onError: (error: Error) => {
      console.error(error);
      options.onError?.(error);
    }
  }));
};

type GenericMutationOptions = {
  onSuccess?: () => void;
  onError?: (error: Error) => void;
};


export const startApplicationMutationQuery = (data: { projectId: string, environmentId: string, resourceID: string }, options: GenericMutationOptions) => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: () => applicationsApi.start({ id: data.resourceID, project_id: data.projectId, environment_id: data.environmentId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: applicationsKeys.detail(data.projectId, data.environmentId, data.resourceID) });
      options.onSuccess?.();
    },
    onError: (error: Error) => {
      console.error(error);
      options.onError?.(error);
    }
  }));
};

export const stopApplicationMutationQuery = (data: { projectId: string, environmentId: string, resourceID: string }, options: GenericMutationOptions) => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: () => applicationsApi.stop({ id: data.resourceID, project_id: data.projectId, environment_id: data.environmentId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: applicationsKeys.detail(data.projectId, data.environmentId, data.resourceID) });
      options.onSuccess?.();
    },
    onError: (error: Error) => {
      console.error(error);
      options.onError?.(error);
    }
  }));
};

export const restartApplicationMutationQuery = (data: { projectId: string, environmentId: string, resourceID: string }, options: GenericMutationOptions) => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: () => applicationsApi.restart({ id: data.resourceID, project_id: data.projectId, environment_id: data.environmentId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: applicationsKeys.detail(data.projectId, data.environmentId, data.resourceID) });
      options.onSuccess?.();
    },
    onError: (error: Error) => {
      console.error(error);
      options.onError?.(error);
    }
  }));
};

export const UpdateApplicationGeneralSettingsMutationQuery = () => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: applicationsApi.updateGeneral,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: applicationsKeys.all });
    },
    onError: (error: Error) => {
      console.error(error);
    }
  }));
};

export const generateApplicationDomainMutationQuery = () => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: ({ projectId, applicationId }: { projectId: string; applicationId: string }) =>
      applicationsApi.generateDomain(projectId, applicationId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: applicationsKeys.all });
    },
    onError: (error: Error) => {
      console.error(error);
    }
  }));
};

export const assignApplicationDomainMutationQuery = () => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: ({
      projectId,
      applicationId,
      domain
    }: {
      projectId: string;
      applicationId: string;
      domain: string;
    }) => applicationsApi.assignDomain(projectId, applicationId, domain),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: applicationsKeys.all });
    },
    onError: (error: Error) => {
      console.error(error);
    }
  }));
};

export const updateApplicationPortsMutationQuery = () => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: ({
      projectId,
      applicationId,
      data
    }: {
      projectId: string;
      applicationId: string;
      data: {
        exposed_ports: number[];
        port_mappings: Array<{ host_port: number; container_port: number }>;
      };
    }) => applicationsApi.updatePorts(projectId, applicationId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: applicationsKeys.all });
    },
    onError: (error: Error) => {
      console.error(error);
    }
  }));
};

export const updateApplicationMutationQuery = () => {
  const queryClient = useQueryClient();

  return createMutation(() => ({
    mutationFn: applicationsApi.updateGeneral,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: applicationsKeys.all });
    },
    onError: (error: Error) => {
      console.error(error);
    }
  }));
};
