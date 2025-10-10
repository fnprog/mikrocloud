import { createQuery } from '@tanstack/svelte-query';
import { templatesApi } from '../api';
import { templatesKeys } from '../keys';

export const createTemplatesListQuery = (category?: string) =>
	createQuery(() => ({
		queryKey: templatesKeys.list(category),
		queryFn: () => templatesApi.list(category)
	}));

export const createTemplateQuery = (templateId: string) =>
	createQuery(() => ({
		queryKey: templatesKeys.detail(templateId),
		queryFn: () => templatesApi.get(templateId)
	}));
