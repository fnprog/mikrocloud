import { apiClient } from '$lib/api/client';
import type { GeneralSettings, UpdateGeneralSettingsRequest, DetectedIPs, SMTPSettings, UpdateSMTPSettingsRequest } from '../types';

export const settingsApi = {
	getGeneral: async (): Promise<GeneralSettings> => {
		return apiClient.get<GeneralSettings>('/settings/general');
	},

	updateGeneral: async (data: UpdateGeneralSettingsRequest): Promise<void> => {
		return apiClient.post<void>('/settings/general', data);
	},

	getSMTP: async (): Promise<SMTPSettings> => {
		return apiClient.get<SMTPSettings>('/settings/smtp');
	},

	updateSMTP: async (data: UpdateSMTPSettingsRequest): Promise<void> => {
		return apiClient.post<void>('/settings/smtp', data);
	},

	detectIPs: async (): Promise<DetectedIPs> => {
		return apiClient.get<DetectedIPs>('/settings/detect-ips');
	}
};
