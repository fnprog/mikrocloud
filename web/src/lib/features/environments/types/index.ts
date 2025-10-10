export interface Environment {
	id: string;
	name: string;
	project_id: string;
	description: string;
	is_production: boolean;
	variables?: Record<string, string>;
	created_at: string;
	updated_at?: string;
}

export interface CreateEnvironmentRequest {
	name: string;
	description?: string;
	is_production?: boolean;
	variables?: Record<string, string>;
}

export interface EnvironmentsResponse {
	environments: Environment[];
}
