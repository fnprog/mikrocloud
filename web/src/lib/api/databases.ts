import { apiClient } from './client';

export type DatabaseType = 'postgresql' | 'mysql' | 'mariadb' | 'redis' | 'keydb' | 'dragonfly' | 'mongodb' | 'clickhouse';
export type DatabaseStatus = 'created' | 'provisioning' | 'running' | 'stopped' | 'failed' | 'deleting';

export interface PostgreSQLConfig {
	version: string;
	database_name: string;
	username: string;
	password: string;
	port: number;
	extensions?: string[];
	environment?: Record<string, string>;
}

export interface MySQLConfig {
	version: string;
	database_name: string;
	username: string;
	password: string;
	root_password: string;
	port: number;
	character_set?: string;
	collation?: string;
	environment?: Record<string, string>;
}

export interface MariaDBConfig {
	version: string;
	database_name: string;
	username: string;
	password: string;
	root_password: string;
	port: number;
	character_set?: string;
	collation?: string;
	environment?: Record<string, string>;
}

export interface RedisConfig {
	version: string;
	password?: string;
	port: number;
	database: number;
	max_memory?: string;
	max_memory_policy?: string;
	persistence?: {
		enabled: boolean;
		type: string;
	};
	environment?: Record<string, string>;
}

export interface KeyDBConfig {
	version: string;
	password?: string;
	port: number;
	database: number;
	max_memory?: string;
	max_memory_policy?: string;
	persistence?: {
		enabled: boolean;
		type: string;
	};
	environment?: Record<string, string>;
}

export interface DragonflyConfig {
	version: string;
	password?: string;
	port: number;
	max_memory?: string;
	persistence: boolean;
	environment?: Record<string, string>;
}

export interface MongoDBConfig {
	version: string;
	database_name: string;
	username: string;
	password: string;
	port: number;
	auth_source?: string;
	replica_set?: string;
	environment?: Record<string, string>;
}

export interface ClickHouseConfig {
	version: string;
	database_name: string;
	username: string;
	password: string;
	port: number;
	http_port?: number;
	environment?: Record<string, string>;
}

export interface DatabaseConfig {
	type: DatabaseType;
	postgresql?: PostgreSQLConfig;
	mysql?: MySQLConfig;
	mariadb?: MariaDBConfig;
	redis?: RedisConfig;
	keydb?: KeyDBConfig;
	dragonfly?: DragonflyConfig;
	mongodb?: MongoDBConfig;
	clickhouse?: ClickHouseConfig;
}

export interface Database {
	id: string;
	name: string;
	description: string;
	type: DatabaseType;
	project_id: string;
	environment_id: string;
	status: DatabaseStatus;
	config: DatabaseConfig;
	connection_string: string;
	ports: Record<string, number>;
	created_at: string;
	updated_at: string;
}

export interface CreateDatabaseRequest {
	name: string;
	description?: string;
	type: DatabaseType;
	environment_id: string;
	config?: DatabaseConfig;
}

export interface DatabasesResponse {
	databases: Database[];
}

export interface DatabaseTypesResponse {
	types: DatabaseType[];
}

export const databasesApi = {
	async list(projectId: string, environmentId?: string): Promise<Database[]> {
		const url = environmentId
			? `/projects/${projectId}/databases?environment_id=${environmentId}`
			: `/projects/${projectId}/databases`;
		const response = await apiClient.get<DatabasesResponse>(url);
		return response.databases;
	},

	async get(projectId: string, databaseId: string): Promise<Database> {
		return apiClient.get<Database>(`/projects/${projectId}/databases/${databaseId}`);
	},

	async create(projectId: string, data: CreateDatabaseRequest): Promise<Database> {
		return apiClient.post<Database>(`/projects/${projectId}/databases`, data);
	},

	async delete(projectId: string, databaseId: string): Promise<void> {
		return apiClient.delete<void>(`/projects/${projectId}/databases/${databaseId}`);
	},

	async getTypes(projectId: string): Promise<DatabaseType[]> {
		const response = await apiClient.get<DatabaseTypesResponse>(`/projects/${projectId}/databases/types`);
		return response.types;
	},

	async getDefaultConfig(projectId: string, type: DatabaseType): Promise<DatabaseConfig> {
		return apiClient.get<DatabaseConfig>(`/projects/${projectId}/databases/types/${type}/config`);
	}
};
