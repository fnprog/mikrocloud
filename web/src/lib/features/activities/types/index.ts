export interface Activity {
  id: string;
  event_type: string;
  level: ActivityLevel;
  description?: string;
  initiator: Initiator;
  metadata?: Record<string, any>;
  created_at: string;
}

export interface Initiator {
  name: string;
  avatar?: string;
  initials: string;
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

export type ActivityLevel = 'info' | 'error' | 'warn' | 'success';
