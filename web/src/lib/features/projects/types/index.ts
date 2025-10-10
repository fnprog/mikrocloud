export interface Project {
  id: string;
  name: string;
  description?: string;
  created_at: string;
}

export interface CreateProjectRequest {
  name: string;
  description?: string;
}

export interface ProjectsResponse {
  projects: Project[];
}
