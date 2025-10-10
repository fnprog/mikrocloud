import { apiClient } from '$lib/api/client';
import type { ActivitiesResponse, ActivitiesForResourceRequest, ActivitiesRequest } from "../types"

//TODO: Fix the endpoint backend wise (org is from the jwt so we should not need it)
export const activitiesApi = {
  async getRecent({ organization_id, limit = 20, offset = 0 }: ActivitiesRequest): Promise<ActivitiesResponse> {
    return apiClient.get<ActivitiesResponse>(`/activities/${organization_id}?limit=${limit}&offset=${offset}`);
  },

  async getForResource({ resource_type, resource_id, limit = 20, offset = 0 }: ActivitiesForResourceRequest): Promise<ActivitiesResponse> {
    return apiClient.get<ActivitiesResponse>(`/activities/${resource_type}/${resource_id}?limit=${limit}&offset=${offset}`);
  },
};
