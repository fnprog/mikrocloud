import { apiClient } from '$lib/api/client';
import type { ServiceTemplate, DeployTemplateRequest, TemplatesResponse } from '../types';

export const templatesApi = {
	list: async (category?: string): Promise<ServiceTemplate[]> => {
		const params = new URLSearchParams();
		if (category) {
			params.append('category', category);
		}
		const response = await apiClient.get<TemplatesResponse>(
			`/templates${params.toString() ? `?${params.toString()}` : ''}`
		);
		return response.templates;
	},

	get: async (id: string): Promise<ServiceTemplate> => {
		return await apiClient.get<ServiceTemplate>(`/templates/${id}`);
	},

	deploy: async (templateId: string, request: DeployTemplateRequest) => {
		return await apiClient.post(`/templates/${templateId}/deploy`, request);
	}
};
