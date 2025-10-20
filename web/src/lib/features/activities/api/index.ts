import { apiClient } from '$lib/api/client';
import type { ActivitiesResponse, ActivitiesForResourceRequest, ActivitiesRequest } from "../types"

export const activitiesApi = {
  async getRecent({ limit = 20, offset = 0 }: ActivitiesRequest): Promise<ActivitiesResponse> {
    return apiClient.get<ActivitiesResponse>(`/activities/?limit=${limit}&offset=${offset}`);
  },

  async getForResource({ resource_type, resource_id, limit = 20, offset = 0 }: ActivitiesForResourceRequest): Promise<ActivitiesResponse> {
    return apiClient.get<ActivitiesResponse>(`/activities/${resource_type}/${resource_id}?limit=${limit}&offset=${offset}`);
  },
};
