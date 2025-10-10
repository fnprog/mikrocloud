import { apiClient } from '$lib/api/client';
import type {
	ValidateRepositoryRequest,
	ValidateRepositoryResponse,
	ListBranchesRequest,
	ListBranchesResponse,
	DetectBuildMethodRequest,
	DetectBuildMethodResponse,
	GitSource,
	GitSourcesResponse,
	CreateGitSourceRequest,
	UpdateGitSourceRequest
} from '../types';

export const gitApi = {
	async validateRepository(
		req: ValidateRepositoryRequest
	): Promise<ValidateRepositoryResponse> {
		return apiClient.post<ValidateRepositoryResponse>('/git/validate', req);
	},

	async listBranches(req: ListBranchesRequest): Promise<ListBranchesResponse> {
		return apiClient.post<ListBranchesResponse>('/git/branches', req);
	},

	async detectBuildMethod(req: DetectBuildMethodRequest): Promise<DetectBuildMethodResponse> {
		return apiClient.post<DetectBuildMethodResponse>('/git/detect-build', req);
	},

	async listGitSources(): Promise<GitSource[]> {
		const response = await apiClient.get<GitSourcesResponse>('/git/sources');
		return response.sources;
	},

	async getGitSource(sourceId: string): Promise<GitSource> {
		return apiClient.get<GitSource>(`/git/sources/${sourceId}`);
	},

	async createGitSource(req: CreateGitSourceRequest): Promise<GitSource> {
		return apiClient.post<GitSource>('/git/sources', req);
	},

	async updateGitSource(sourceId: string, req: UpdateGitSourceRequest): Promise<GitSource> {
		return apiClient.put<GitSource>(`/git/sources/${sourceId}`, req);
	},

	async deleteGitSource(sourceId: string): Promise<void> {
		return apiClient.delete<void>(`/git/sources/${sourceId}`);
	}
};
