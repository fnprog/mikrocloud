export type ApplicationStatus = 'pending' | 'running' | 'stopped' | 'failed' | 'deploying';

export interface Application {
  id: string;
  name: string;
  description: string;
  project_id: string;
  environment_id: string;
  domain: string;
  status: ApplicationStatus;
  created_at: string;
  custom_domain?: string;
  generated_domain?: string;
  exposed_ports?: number[];
  port_mappings?: Array<{ host: number; container: number }>;
}

export interface GitRepoSource {
  url: string;
  branch: string;
  path?: string;
  base_path?: string;
  token?: string;
}

export interface RegistrySource {
  image: string;
  tag: string;
}

export interface UploadSource {
  filename: string;
  file_path: string;
}

export interface DeploymentSource {
  type: 'git' | 'registry' | 'upload';
  git_repo?: GitRepoSource;
  registry?: RegistrySource;
  upload?: UploadSource;
}

export interface CreateApplicationRequest {
  name: string;
  description?: string;
  project_id: string;
  environment_id: string;
  deployment_source: DeploymentSource;
  buildpack: string;
  env_vars?: Record<string, string>;
}


export interface DeleteApplicationRequest {
  id: string;
  project_id: string;
  environment_id: string;
}


export interface SwitchApplicationStateRequest {
  id: string;
  project_id: string;
  environment_id: string;
}

export interface ApplicationsResponse {
  applications: Application[];
}

export interface UpdateGeneralSettingsRequest {
  project_id: string,
  application_id: string,
  payload: { name?: string; description?: string }
}
