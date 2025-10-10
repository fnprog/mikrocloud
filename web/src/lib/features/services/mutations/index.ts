import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { templatesApi } from '../api';
import { applicationsKeys } from '../../applications/keys';
import type { DeployTemplateRequest } from '../types';

export const deployTemplateMutationQuery = (templateId: string) => {
	const queryClient = useQueryClient();
	return createMutation(() => ({
		mutationFn: (data: DeployTemplateRequest) => templatesApi.deploy(templateId, data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: applicationsKeys.lists() });
		}
	}));
};
