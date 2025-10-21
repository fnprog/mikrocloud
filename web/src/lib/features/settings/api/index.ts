import { apiClient } from '$lib/api/client';
import type { GeneralSettings, UpdateGeneralSettingsRequest, DetectedIPs } from '../types';

export const settingsApi = {
  getGeneral: async (): Promise<GeneralSettings> => {
    return apiClient.get<GeneralSettings>('/settings/general');
  },

  updateGeneral: async (data: UpdateGeneralSettingsRequest): Promise<void> => {
    return apiClient.post<void>('/settings/general', data);
  },

  detectIPs: async (): Promise<DetectedIPs> => {
    return apiClient.get<DetectedIPs>('/settings/detect-ips');
  }
};
