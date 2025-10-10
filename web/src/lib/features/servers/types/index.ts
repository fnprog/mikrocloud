export interface Server {
	id: string;
	name: string;
	description?: string;
	hostname: string;
	ip_address: string;
	port: number;
	ssh_key?: string;
	server_type: 'control_plane' | 'worker' | 'database' | 'proxy';
	status: 'online' | 'offline' | 'maintenance' | 'error' | 'unknown';
	cpu_cores?: number;
	memory_mb?: number;
	disk_gb?: number;
	os?: string;
	os_version?: string;
	metadata?: string;
	tags?: string[];
	organization_id: string;
	created_at: string;
	updated_at: string;
}

export interface CreateServerRequest {
	name: string;
	description?: string;
	hostname: string;
	ip_address: string;
	port: number;
	ssh_key?: string;
	server_type: 'control_plane' | 'worker' | 'database' | 'proxy';
	tags?: string[];
}

export interface UpdateServerRequest {
	name?: string;
	description?: string;
	hostname?: string;
	ip_address?: string;
	port?: number;
	ssh_key?: string;
	status?: 'online' | 'offline' | 'maintenance' | 'error' | 'unknown';
	cpu_cores?: number;
	memory_mb?: number;
	disk_gb?: number;
	os?: string;
	os_version?: string;
	tags?: string[];
}
