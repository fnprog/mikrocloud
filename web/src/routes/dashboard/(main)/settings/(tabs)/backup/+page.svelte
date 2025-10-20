<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import * as Card from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import { Download, Trash2 } from 'lucide-svelte';

	let dbType = $state('sqlite');
	let dbUuid = $state('');
	let dbName = $state('');
	let dbDescription = $state('');
	let dbUser = $state('');
	let dbPassword = $state('');

	let backupEnabled = $state(true);
	let s3BackupEnabled = $state(false);
	let disableLocalBackup = $state(false);
	let cronFrequency = $state('0 0 * * *');
	let timeout = $state(300);

	let retentionCount = $state(7);
	let retentionDays = $state(30);
	let maxStorageGb = $state(0);

	type BackupTask = {
		id: string;
		timestamp: string;
		status: 'success' | 'failed' | 'deleted';
		size: string;
		location: string;
	};

	let tasks = $state<BackupTask[]>([]);

	onMount(async () => {
		try {
			const response = await fetch('/api/settings/backup');
			if (response.ok) {
				const data = await response.json();
				dbType = data.db_type || 'sqlite';
				dbUuid = data.db_uuid || '';
				dbName = data.db_name || '';
				dbDescription = data.db_description || '';
				dbUser = data.db_user || '';
				dbPassword = data.db_password || '';

				backupEnabled = data.backup_enabled ?? true;
				s3BackupEnabled = data.s3_backup_enabled ?? false;
				disableLocalBackup = data.disable_local_backup ?? false;
				cronFrequency = data.cron_frequency || '0 0 * * *';
				timeout = data.timeout || 300;

				retentionCount = data.retention_count || 7;
				retentionDays = data.retention_days || 30;
				maxStorageGb = data.max_storage_gb || 0;

				tasks = data.tasks || [];
			}
		} catch (error) {
			console.error('Failed to load backup settings:', error);
		}
	});

	async function handleSave() {
		try {
			const response = await fetch('/api/settings/backup', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					backup_enabled: backupEnabled,
					s3_backup_enabled: s3BackupEnabled,
					disable_local_backup: disableLocalBackup,
					cron_frequency: cronFrequency,
					timeout,
					retention_count: retentionCount,
					retention_days: retentionDays,
					max_storage_gb: maxStorageGb
				})
			});

			if (!response.ok) {
				throw new Error('Failed to save backup settings');
			}
		} catch (error) {
			console.error('Failed to save backup settings:', error);
		}
	}

	async function handleDownload(taskId: string) {
		try {
			const response = await fetch(`/api/settings/backup/download/${taskId}`);
			if (!response.ok) throw new Error('Download failed');

			const blob = await response.blob();
			const url = window.URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = `backup-${taskId}.tar.gz`;
			document.body.appendChild(a);
			a.click();
			window.URL.revokeObjectURL(url);
			document.body.removeChild(a);
		} catch (error) {
			console.error('Failed to download backup:', error);
		}
	}

	async function handleDelete(taskId: string) {
		try {
			const response = await fetch(`/api/settings/backup/${taskId}`, { method: 'DELETE' });
			if (!response.ok) throw new Error('Delete failed');
			tasks = tasks.filter((t) => t.id !== taskId);
		} catch (error) {
			console.error('Failed to delete backup:', error);
		}
	}

	async function handleCleanupFailed() {
		try {
			const response = await fetch('/api/settings/backup/cleanup/failed', { method: 'DELETE' });
			if (!response.ok) throw new Error('Cleanup failed');
			tasks = tasks.filter((t) => t.status !== 'failed');
		} catch (error) {
			console.error('Failed to cleanup failed backups:', error);
		}
	}

	async function handleCleanupDeleted() {
		try {
			const response = await fetch('/api/settings/backup/cleanup/deleted', { method: 'DELETE' });
			if (!response.ok) throw new Error('Cleanup failed');
			tasks = tasks.filter((t) => t.status !== 'deleted');
		} catch (error) {
			console.error('Failed to cleanup deleted backups:', error);
		}
	}
</script>

<div class="space-y-6">
	<Card.Root>
		<Card.Header>
			<Card.Title>Instance Backup Configuration</Card.Title>
			<Card.Description>Database information for backup operations.</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label>Database Type</Label>
					<Input value={dbType} readonly class="bg-muted" />
				</div>
				<div class="space-y-2">
					<Label>Database UUID</Label>
					<Input value={dbUuid} readonly class="bg-muted" />
				</div>
			</div>
			<div class="space-y-2">
				<Label>Database Name</Label>
				<Input value={dbName} readonly class="bg-muted" />
			</div>
			<div class="space-y-2">
				<Label>Description</Label>
				<Input value={dbDescription} readonly class="bg-muted" />
			</div>
			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label>User</Label>
					<Input value={dbUser} readonly class="bg-muted" />
				</div>
				<div class="space-y-2">
					<Label>Password</Label>
					<Input type="password" value={dbPassword} readonly class="bg-muted" />
				</div>
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title>Scheduled Backup</Card.Title>
			<Card.Description>Configure automatic backup schedules and destinations.</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-6">
			<div class="flex items-center justify-between">
				<div class="space-y-1">
					<Label for="backup-enabled">Backup Enabled</Label>
					<p class="text-xs text-muted-foreground">Enable automatic backups.</p>
				</div>
				<Switch id="backup-enabled" bind:checked={backupEnabled} />
			</div>

			<div class="flex items-center justify-between">
				<div class="space-y-1">
					<Label for="s3-backup">S3 Backup Enabled</Label>
					<p class="text-xs text-muted-foreground">Upload backups to S3-compatible storage.</p>
				</div>
				<Switch id="s3-backup" bind:checked={s3BackupEnabled} />
			</div>

			<div class="flex items-center justify-between">
				<div class="space-y-1">
					<Label for="disable-local">Disable Local Backup</Label>
					<p class="text-xs text-muted-foreground">Only store backups in S3, not locally.</p>
				</div>
				<Switch id="disable-local" bind:checked={disableLocalBackup} />
			</div>

			<div class="space-y-2">
				<Label for="cron">Cron Frequency</Label>
				<Input id="cron" bind:value={cronFrequency} placeholder="0 0 * * *" />
				<p class="text-xs text-muted-foreground">
					Cron expression for backup schedule. Default is <code class="rounded bg-muted px-1 py-0.5"
						>0 0 * * *</code
					> (daily at midnight).
				</p>
			</div>

			<div class="space-y-2">
				<Label for="timeout">Timeout (seconds)</Label>
				<Input id="timeout" type="number" bind:value={timeout} />
				<p class="text-xs text-muted-foreground">
					Maximum time allowed for backup operations. Default is 300 seconds.
				</p>
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title>Backup Retention Settings</Card.Title>
			<Card.Description>Configure how long backups are kept.</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<div class="space-y-2">
				<Label for="retention-count">Number of Backups to Keep</Label>
				<Input id="retention-count" type="number" bind:value={retentionCount} />
				<p class="text-xs text-muted-foreground">
					Maximum number of backups to retain. Older backups will be deleted.
				</p>
			</div>

			<div class="space-y-2">
				<Label for="retention-days">Days to Keep Backups</Label>
				<Input id="retention-days" type="number" bind:value={retentionDays} />
				<p class="text-xs text-muted-foreground">
					Number of days to keep backups. Backups older than this will be deleted.
				</p>
			</div>

			<div class="space-y-2">
				<Label for="max-storage">Maximum Storage (GB)</Label>
				<Input id="max-storage" type="number" bind:value={maxStorageGb} />
				<p class="text-xs text-muted-foreground">
					Maximum storage for backups in GB. Set to 0 for unlimited.
				</p>
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<div class="flex items-center justify-between">
				<div>
					<Card.Title>Tasks Execution List</Card.Title>
					<Card.Description>Recent backup tasks and their status.</Card.Description>
				</div>
				<div class="flex gap-2">
					<Button variant="outline" size="sm" onclick={handleCleanupFailed}>
						Cleanup Failed
					</Button>
					<Button variant="outline" size="sm" onclick={handleCleanupDeleted}>
						Cleanup Deleted
					</Button>
				</div>
			</div>
		</Card.Header>
		<Card.Content>
			{#if tasks.length === 0}
				<p class="text-sm text-muted-foreground text-center py-8">No backup tasks yet.</p>
			{:else}
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>Timestamp</Table.Head>
							<Table.Head>Status</Table.Head>
							<Table.Head>Size</Table.Head>
							<Table.Head>Location</Table.Head>
							<Table.Head class="text-right">Actions</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each tasks as task}
							<Table.Row>
								<Table.Cell class="font-medium">{task.timestamp}</Table.Cell>
								<Table.Cell>
									<span
										class="inline-flex items-center rounded-full px-2 py-1 text-xs font-medium"
										class:bg-green-50={task.status === 'success'}
										class:text-green-700={task.status === 'success'}
										class:dark:bg-green-950={task.status === 'success'}
										class:dark:text-green-400={task.status === 'success'}
										class:bg-red-50={task.status === 'failed'}
										class:text-red-700={task.status === 'failed'}
										class:dark:bg-red-950={task.status === 'failed'}
										class:dark:text-red-400={task.status === 'failed'}
										class:bg-gray-50={task.status === 'deleted'}
										class:text-gray-700={task.status === 'deleted'}
										class:dark:bg-gray-950={task.status === 'deleted'}
										class:dark:text-gray-400={task.status === 'deleted'}
									>
										{task.status}
									</span>
								</Table.Cell>
								<Table.Cell>{task.size}</Table.Cell>
								<Table.Cell>{task.location}</Table.Cell>
								<Table.Cell class="text-right">
									<div class="flex justify-end gap-2">
										{#if task.status === 'success'}
											<Button
												variant="ghost"
												size="sm"
												onclick={() => handleDownload(task.id)}
											>
												<Download class="h-4 w-4" />
											</Button>
										{/if}
										<Button
											variant="ghost"
											size="sm"
											onclick={() => handleDelete(task.id)}
										>
											<Trash2 class="h-4 w-4" />
										</Button>
									</div>
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			{/if}
		</Card.Content>
	</Card.Root>

	<div class="flex justify-end">
		<Button onclick={handleSave}>Save Changes</Button>
	</div>
</div>
