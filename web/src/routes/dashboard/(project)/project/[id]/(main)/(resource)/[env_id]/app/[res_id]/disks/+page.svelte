<script>
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import {
		Plus,
		HardDrive,
		Archive,
		Trash2,
		Download,
		RefreshCw,
		AlertTriangle,
		CheckCircle,
		Clock
	} from 'lucide-svelte';

	const projectId = page.params.id;

	let project = $state({
		name: 'Focalpoint Dashboard',
		workspace: 'Focalpoint',
		category: 'Applications'
	});

	let showCreateDiskModal = $state(false);
	let showBackupModal = $state(false);
	let selectedDisk = $state(null);

	let newDisk = $state({
		name: '',
		size: 10,
		type: 'ssd',
		encrypted: true
	});

	let newBackup = $state({
		name: '',
		description: ''
	});

	// Mock disk data
	let disks = $state([
		{
			id: 1,
			name: 'root-disk',
			type: 'ssd',
			size: 50,
			used: 32,
			available: 18,
			mountPoint: '/',
			status: 'healthy',
			encrypted: true,
			createdAt: '2024-01-15',
			lastBackup: '2024-01-30'
		},
		{
			id: 2,
			name: 'data-disk',
			type: 'ssd',
			size: 100,
			used: 67,
			available: 33,
			mountPoint: '/data',
			status: 'healthy',
			encrypted: true,
			createdAt: '2024-01-20',
			lastBackup: '2024-01-29'
		},
		{
			id: 3,
			name: 'logs-disk',
			type: 'hdd',
			size: 200,
			used: 145,
			available: 55,
			mountPoint: '/var/log',
			status: 'warning',
			encrypted: false,
			createdAt: '2024-01-25',
			lastBackup: '2024-01-28'
		}
	]);

	let backups = $state([
		{
			id: 1,
			diskId: 1,
			name: 'root-disk-backup-20240130',
			size: 28,
			status: 'completed',
			createdAt: '2024-01-30 14:30:00',
			type: 'automatic'
		},
		{
			id: 2,
			diskId: 2,
			name: 'data-disk-backup-20240129',
			size: 62,
			status: 'completed',
			createdAt: '2024-01-29 02:00:00',
			type: 'automatic'
		},
		{
			id: 3,
			diskId: 1,
			name: 'pre-deployment-backup',
			size: 27,
			status: 'completed',
			createdAt: '2024-01-28 10:15:00',
			type: 'manual'
		}
	]);

	function getStatusColor(status) {
		switch (status) {
			case 'healthy':
				return 'text-green-600 bg-green-50';
			case 'warning':
				return 'text-yellow-600 bg-yellow-50';
			case 'error':
				return 'text-red-600 bg-red-50';
			default:
				return 'text-gray-600 bg-gray-50';
		}
	}

	function getStatusIcon(status) {
		switch (status) {
			case 'healthy':
				return CheckCircle;
			case 'warning':
				return AlertTriangle;
			case 'error':
				return AlertTriangle;
			default:
				return Clock;
		}
	}

	function getUsageColor(percentage) {
		if (percentage > 90) return 'bg-red-500';
		if (percentage > 75) return 'bg-yellow-500';
		return 'bg-green-500';
	}

	function createDisk() {
		if (newDisk.name && newDisk.size) {
			const newId = Math.max(...disks.map((d) => d.id)) + 1;
			disks = [
				...disks,
				{
					id: newId,
					name: newDisk.name,
					type: newDisk.type,
					size: newDisk.size,
					used: 0,
					available: newDisk.size,
					mountPoint: `/mnt/${newDisk.name}`,
					status: 'healthy',
					encrypted: newDisk.encrypted,
					createdAt: new Date().toISOString().split('T')[0],
					lastBackup: null
				}
			];
			newDisk = { name: '', size: 10, type: 'ssd', encrypted: true };
			showCreateDiskModal = false;
		}
	}

	function createBackup(diskId) {
		selectedDisk = disks.find((d) => d.id === diskId);
		showBackupModal = true;
	}

	function executeBackup() {
		if (selectedDisk && newBackup.name) {
			const newId = Math.max(...backups.map((b) => b.id)) + 1;
			const disk = disks.find((d) => d.id === selectedDisk.id);
			backups = [
				...backups,
				{
					id: newId,
					diskId: selectedDisk.id,
					name: newBackup.name,
					size: disk.used,
					status: 'completed',
					createdAt: new Date().toISOString().replace('T', ' ').split('.')[0],
					type: 'manual'
				}
			];
			disk.lastBackup = new Date().toISOString().split('T')[0];
			newBackup = { name: '', description: '' };
			showBackupModal = false;
			selectedDisk = null;
		}
	}

	function deleteDisk(diskId) {
		disks = disks.filter((d) => d.id !== diskId);
		backups = backups.filter((b) => b.diskId !== diskId);
	}

	function deleteBackup(backupId) {
		backups = backups.filter((b) => b.id !== backupId);
	}

	function formatBytes(bytes) {
		return `${bytes} GB`;
	}

	function getUsagePercentage(used, total) {
		return Math.round((used / total) * 100);
	}
</script>

<svelte:head>
	<title>Disks - {project.name}</title>
</svelte:head>

<div class="flex-1 p-6">
	<!-- Header -->
	<div class="flex items-center justify-between mb-6">
		<div class="flex items-center space-x-4">
			<div>
				<h1 class="text-2xl font-semibold text-gray-900">Disks</h1>
				<p class="text-sm text-gray-500 mt-1">
					Manage storage volumes and backups for {project.name}.
				</p>
			</div>
		</div>
		<Button onclick={() => (showCreateDiskModal = true)}>
			<Plus class="w-4 h-4 mr-2" />
			Create Disk
		</Button>
	</div>

	<!-- Disks List -->
	<div class="space-y-6 mb-8">
		{#each disks as disk}
			<Card>
				<CardContent class="p-6">
					<div class="flex items-center justify-between mb-4">
						<div class="flex items-center space-x-3">
							<HardDrive class="w-6 h-6 text-gray-600" />
							<div>
								<h3 class="font-medium text-gray-900">{disk.name}</h3>
								<p class="text-sm text-gray-600">{disk.mountPoint}</p>
							</div>
							<Badge variant="outline" class="text-xs {getStatusColor(disk.status)}">
								{disk.status}
							</Badge>
							{#if disk.encrypted}
								<Badge variant="outline" class="text-xs">Encrypted</Badge>
							{/if}
						</div>
						<div class="flex items-center space-x-2">
							<Button size="sm" variant="outline" onclick={() => createBackup(disk.id)}>
								<Archive class="w-4 h-4 mr-1" />
								Backup
							</Button>
							<Button size="sm" variant="outline">
								<RefreshCw class="w-4 h-4" />
							</Button>
							<Button size="sm" variant="outline" onclick={() => deleteDisk(disk.id)}>
								<Trash2 class="w-4 h-4" />
							</Button>
						</div>
					</div>

					<div class="grid grid-cols-1 md:grid-cols-4 gap-6">
						<div>
							<p class="text-sm text-gray-600">Type</p>
							<p class="font-medium text-gray-900">{disk.type.toUpperCase()}</p>
						</div>
						<div>
							<p class="text-sm text-gray-600">Total Size</p>
							<p class="font-medium text-gray-900">{formatBytes(disk.size)}</p>
						</div>
						<div>
							<p class="text-sm text-gray-600">Used</p>
							<p class="font-medium text-gray-900">
								{formatBytes(disk.used)} ({getUsagePercentage(disk.used, disk.size)}%)
							</p>
						</div>
						<div>
							<p class="text-sm text-gray-600">Last Backup</p>
							<p class="font-medium text-gray-900">{disk.lastBackup || 'Never'}</p>
						</div>
					</div>

					<!-- Usage Bar -->
					<div class="mt-4">
						<div class="flex items-center justify-between text-sm text-gray-600 mb-1">
							<span>Storage Usage</span>
							<span>{getUsagePercentage(disk.used, disk.size)}% used</span>
						</div>
						<div class="w-full bg-gray-200 rounded-full h-2">
							<div
								class="h-2 rounded-full {getUsageColor(getUsagePercentage(disk.used, disk.size))}"
								style="width: {getUsagePercentage(disk.used, disk.size)}%"
							></div>
						</div>
					</div>
				</CardContent>
			</Card>
		{/each}
	</div>

	<!-- Backups -->
	<Card>
		<CardHeader>
			<CardTitle>Recent Backups</CardTitle>
		</CardHeader>
		<CardContent>
			<div class="space-y-4">
				{#each backups as backup}
					<div class="flex items-center justify-between p-4 border rounded-lg">
						<div class="flex items-center space-x-4">
							<Archive class="w-5 h-5 text-gray-600" />
							<div>
								<h4 class="font-medium text-gray-900">{backup.name}</h4>
								<div class="flex items-center space-x-4 text-sm text-gray-600">
									<span>Size: {formatBytes(backup.size)}</span>
									<span>•</span>
									<span>Created: {backup.createdAt}</span>
									<span>•</span>
									<Badge variant="outline" class="text-xs">
										{backup.type}
									</Badge>
								</div>
							</div>
						</div>
						<div class="flex items-center space-x-2">
							<Button size="sm" variant="outline">
								<Download class="w-4 h-4" />
							</Button>
							<Button size="sm" variant="outline" onclick={() => deleteBackup(backup.id)}>
								<Trash2 class="w-4 h-4" />
							</Button>
						</div>
					</div>
				{/each}
			</div>
		</CardContent>
	</Card>
</div>
<!-- Create Disk Modal -->
{#if showCreateDiskModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<Card class="w-full max-w-md">
			<CardHeader>
				<CardTitle>Create New Disk</CardTitle>
			</CardHeader>
			<CardContent class="space-y-4">
				<div>
					<Label for="disk-name">Disk Name</Label>
					<Input id="disk-name" bind:value={newDisk.name} placeholder="my-disk" />
				</div>
				<div>
					<Label for="disk-size">Size (GB)</Label>
					<Input id="disk-size" type="number" bind:value={newDisk.size} min="1" max="1000" />
				</div>
				<div>
					<Label for="disk-type">Type</Label>
					<select
						id="disk-type"
						bind:value={newDisk.type}
						class="w-full px-3 py-2 border border-gray-300 rounded-md"
					>
						<option value="ssd">SSD (High Performance)</option>
						<option value="hdd">HDD (Standard)</option>
					</select>
				</div>
				<div class="flex items-center space-x-2">
					<input type="checkbox" id="encrypted" bind:checked={newDisk.encrypted} class="rounded" />
					<Label for="encrypted">Enable encryption</Label>
				</div>
				<div class="flex space-x-2 pt-4">
					<Button onclick={createDisk} class="flex-1">Create Disk</Button>
					<Button variant="outline" onclick={() => (showCreateDiskModal = false)} class="flex-1"
						>Cancel</Button
					>
				</div>
			</CardContent>
		</Card>
	</div>
{/if}

<!-- Backup Modal -->
{#if showBackupModal && selectedDisk}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<Card class="w-full max-w-md">
			<CardHeader>
				<CardTitle>Create Backup - {selectedDisk.name}</CardTitle>
			</CardHeader>
			<CardContent class="space-y-4">
				<div>
					<Label for="backup-name">Backup Name</Label>
					<Input
						id="backup-name"
						bind:value={newBackup.name}
						placeholder="{selectedDisk.name}-backup-{new Date().toISOString().split('T')[0]}"
					/>
				</div>
				<div>
					<Label for="backup-description">Description (optional)</Label>
					<Input
						id="backup-description"
						bind:value={newBackup.description}
						placeholder="Manual backup before deployment"
					/>
				</div>
				<div class="bg-gray-50 rounded-lg p-4">
					<h4 class="font-medium text-gray-900 mb-2">Backup Details</h4>
					<div class="space-y-1 text-sm text-gray-600">
						<p>Disk: {selectedDisk.name}</p>
						<p>Size: ~{formatBytes(selectedDisk.used)}</p>
						<p>Type: Manual backup</p>
					</div>
				</div>
				<div class="flex space-x-2 pt-4">
					<Button onclick={executeBackup} class="flex-1">Create Backup</Button>
					<Button
						variant="outline"
						onclick={() => {
							showBackupModal = false;
							selectedDisk = null;
						}}
						class="flex-1">Cancel</Button
					>
				</div>
			</CardContent>
		</Card>
	</div>
{/if}
