export interface HealthCheckResponse {
	status: string;
	service: string;
	version: string;
	timestamp: string;
}

export interface SystemComponents {
	database: string;
	storage: string;
	container: string;
}

export interface SystemStatusResponse {
	status: string;
	components: SystemComponents;
	timestamp: string;
}

export interface SystemInfoResponse {
	version: string;
	platform: string;
	go_version: string;
	build_time?: string;
	environment: string;
}

export interface ResourcesResponse {
	projects: number;
	applications: number;
	databases: number;
	services: number;
}

export interface DomainInfo {
	id: string;
	name: string;
	verified: boolean;
	ssl_enabled: boolean;
	ssl_expiry?: string;
	service_id?: string;
	service_name?: string;
	created_at: string;
}

export interface DomainListResponse {
	domains: DomainInfo[];
	total: number;
}

export interface AddDomainRequest {
	name: string;
	service_id?: string;
}

export interface EnableSSLRequest {
	provider: string;
	email?: string;
}
