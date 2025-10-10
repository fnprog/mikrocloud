import { apiClient } from '$lib/api/client';
import type { GeneralSettings, UpdateGeneralSettingsRequest, DetectedIPs } from '../types';

export const settingsApi = {
	getGeneral: async (): Promise<GeneralSettings> => {
		return apiClient.get<GeneralSettings>('/api/settings/general');
	},

	updateGeneral: async (data: UpdateGeneralSettingsRequest): Promise<void> => {
		return apiClient.post<void>('/api/settings/general', data);
	},

	detectIPs: async (): Promise<DetectedIPs> => {
		return apiClient.get<DetectedIPs>('/api/settings/detect-ips');
	}
};
