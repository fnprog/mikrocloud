export interface Disk {
	id: string;
	name: string;
	project_id: string;
	service_id?: string;
	size: number;
	size_gb: number;
	mount_path: string;
	filesystem: string;
	status: string;
	persistent: boolean;
	backup_enabled: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreateDiskRequest {
	name: string;
	size_gb: number;
	mount_path: string;
	filesystem: 'ext4' | 'xfs' | 'btrfs' | 'zfs';
	persistent: boolean;
}

export interface ResizeDiskRequest {
	size_gb: number;
}

export interface AttachDiskRequest {
	service_id: string;
}

export interface ListDisksResponse {
	disks: Disk[];
}
