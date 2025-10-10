export type GitProvider = 'github' | 'gitlab' | 'bitbucket' | 'custom';

export interface ValidateRepositoryRequest {
  provider: GitProvider;
  repository: string;
  branch?: string;
  token?: string;
  custom_url?: string;
}

export interface ValidateRepositoryResponse {
  valid: boolean;
  message?: string;
}

export interface ListBranchesRequest {
  provider: GitProvider;
  repository: string;
  token?: string;
  custom_url?: string;
}

export interface Branch {
  name: string;
  protected: boolean;
}

export interface ListBranchesResponse {
  branches: Branch[];
}

export interface DetectBuildMethodRequest {
  provider: GitProvider;
  repository_url: string;
  branch: string;
  token?: string;
  custom_url?: string;
}

export interface DetectBuildMethodResponse {
  build_method: string;
  dockerfile_path?: string;
  compose_path?: string;
  message?: string;
}

export interface CreateGitSourceRequest {
  name: string;
  provider: GitProvider;
  access_token: string;
  refresh_token?: string;
  token_expires_at?: string;
  custom_url?: string;
}

export interface UpdateGitSourceRequest {
  name?: string;
  access_token?: string;
  refresh_token?: string;
  token_expires_at?: string;
  custom_url?: string;
}

export interface GitSource {
  id: string;
  org_id: string;
  user_id: string;
  provider: GitProvider;
  name: string;
  custom_url?: string;
  created_at: string;
  updated_at: string;
}

export interface GitSourcesResponse {
  sources: GitSource[];
}
