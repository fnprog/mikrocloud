export interface Activity {
  id: string;
  organization_id: string;
  activity_type: string;
  description: string;
  initiator_id?: string;
  initiator_name: string;
  resource_type?: string;
  resource_id?: string;
  resource_name?: string;
  metadata?: Record<string, any>;
  created_at: string;
}

export interface ActivitiesRequest {
  limit: number;
  offset: number;
}

export interface ActivitiesForResourceRequest extends ActivitiesRequest {
  resource_type: string;
  resource_id: string;
}

export interface ActivitiesResponse {
  activities: Activity[];
  total: number;
}
