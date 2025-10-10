export interface ServiceTemplate {
	id: string;
	name: string;
	description: string;
	category: string;
	version: string;
	git_url?: {
		url: string;
		branch: string;
		context_root: string;
	};
	environment: Record<string, string>;
	ports: Array<{
		port: number;
		protocol: string;
		public: boolean;
	}>;
	volumes: Array<{
		name: string;
		mount_path: string;
		size?: string;
		read_only: boolean;
	}>;
	official: boolean;
	created_at?: string;
	updated_at?: string;
}

export interface DeployTemplateRequest {
	name: string;
	project_id: string;
	environment_id: string;
	environment?: Record<string, string>;
	custom_config?: Record<string, unknown>;
}

export interface TemplatesResponse {
	templates: ServiceTemplate[];
}
