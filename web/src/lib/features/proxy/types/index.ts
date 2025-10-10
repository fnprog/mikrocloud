export interface TLSConfig {
	enabled: boolean;
	cert_file?: string;
	key_file?: string;
	auto_generated?: boolean;
	domains?: string[];
}

export interface MiddlewareConfig {
	name: string;
	type: string;
	config: Record<string, unknown>;
}

export interface HealthCheckConfig {
	enabled: boolean;
	path: string;
	interval: string;
	timeout: string;
	retries: number;
}

export interface ProxyConfig {
	id: string;
	name: string;
	project_id: string;
	service_name: string;
	container_id: string;
	hostnames: string[];
	target_url: string;
	port: number;
	protocol: 'http' | 'https' | 'tcp' | 'udp';
	path_prefix?: string;
	strip_prefix?: boolean;
	tls?: TLSConfig;
	middlewares: MiddlewareConfig[];
	health_check?: HealthCheckConfig;
	status: string;
	router_name: string;
	traefik_service_name: string;
	rule_host: string;
	created_at: string;
	updated_at: string;
}

export interface CreateProxyConfigRequest {
	name: string;
	service_name: string;
	container_id: string;
	hostnames: string[];
	target_url: string;
	port: number;
	protocol: 'http' | 'https' | 'tcp' | 'udp';
	path_prefix?: string;
	strip_prefix?: boolean;
	tls?: TLSConfig;
	middlewares?: MiddlewareConfig[];
	health_check?: HealthCheckConfig;
}

export interface UpdateProxyConfigRequest {
	name?: string;
	service_name?: string;
	container_id?: string;
	hostnames?: string[];
	target_url?: string;
	port?: number;
	protocol?: 'http' | 'https' | 'tcp' | 'udp';
	path_prefix?: string;
	strip_prefix?: boolean;
	tls?: TLSConfig;
	middlewares?: MiddlewareConfig[];
	health_check?: HealthCheckConfig;
}

export interface ProxyConfigListResponse {
	proxy_configs: ProxyConfig[];
}
