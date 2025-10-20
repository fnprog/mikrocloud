<script lang="ts">
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import {
		Search,
		RefreshCw,
		User,
		Calendar,
		Shield,
		Database,
		Settings,
		GitBranch,
		Trash2,
		FileEdit,
		Lock,
		Unlock
	} from 'lucide-svelte';

	let searchQuery = $state('');
	let actionFilter = $state('all');
	let userFilter = $state('all');
	let isRefreshing = $state(false);

	const actionFilters = [
		{ value: 'all', label: 'All Actions' },
		{ value: 'create', label: 'Create' },
		{ value: 'update', label: 'Update' },
		{ value: 'delete', label: 'Delete' },
		{ value: 'deploy', label: 'Deploy' },
		{ value: 'login', label: 'Login' },
		{ value: 'permission', label: 'Permission' }
	];

	const userFilters = [
		{ value: 'all', label: 'All Users' },
		{ value: 'john.doe', label: 'John Doe' },
		{ value: 'jane.smith', label: 'Jane Smith' },
		{ value: 'bob.johnson', label: 'Bob Johnson' },
		{ value: 'system', label: 'System' }
	];

	const activities = [
		{
			id: 'act_001',
			timestamp: '2024-10-10T14:35:22Z',
			user: 'john.doe',
			displayName: 'John Doe',
			action: 'deploy',
			resource: 'Application',
			resourceName: 'mikrocloud-app',
			environment: 'production',
			details: 'Deployed version v2.4.1 to production',
			status: 'success',
			ipAddress: '192.168.1.100',
			userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)'
		},
		{
			id: 'act_002',
			timestamp: '2024-10-10T14:28:15Z',
			user: 'jane.smith',
			displayName: 'Jane Smith',
			action: 'update',
			resource: 'Database',
			resourceName: 'main-postgres',
			environment: 'production',
			details: 'Updated connection pool settings (max_connections: 100)',
			status: 'success',
			ipAddress: '192.168.1.105',
			userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'
		},
		{
			id: 'act_003',
			timestamp: '2024-10-10T14:20:45Z',
			user: 'bob.johnson',
			displayName: 'Bob Johnson',
			action: 'delete',
			resource: 'Disk',
			resourceName: 'legacy-data-vol',
			environment: 'staging',
			details: 'Deleted unused disk volume (250GB)',
			status: 'success',
			ipAddress: '192.168.1.112',
			userAgent: 'Mozilla/5.0 (X11; Linux x86_64)'
		},
		{
			id: 'act_004',
			timestamp: '2024-10-10T14:15:30Z',
			user: 'system',
			displayName: 'System',
			action: 'create',
			resource: 'Backup',
			resourceName: 'daily-backup-2024-10-10',
			environment: 'production',
			details: 'Automated daily backup completed (2.4GB)',
			status: 'success',
			ipAddress: 'internal',
			userAgent: 'mikrocloud-scheduler/1.0'
		},
		{
			id: 'act_005',
			timestamp: '2024-10-10T14:12:00Z',
			user: 'jane.smith',
			displayName: 'Jane Smith',
			action: 'permission',
			resource: 'Project',
			resourceName: 'web-dashboard',
			environment: null,
			details: 'Granted deploy permission to alice@example.com',
			status: 'success',
			ipAddress: '192.168.1.105',
			userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'
		},
		{
			id: 'act_006',
			timestamp: '2024-10-10T14:05:20Z',
			user: 'john.doe',
			displayName: 'John Doe',
			action: 'login',
			resource: 'Auth',
			resourceName: null,
			environment: null,
			details: 'User logged in successfully',
			status: 'success',
			ipAddress: '192.168.1.100',
			userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)'
		},
		{
			id: 'act_007',
			timestamp: '2024-10-10T14:00:12Z',
			user: 'bob.johnson',
			displayName: 'Bob Johnson',
			action: 'create',
			resource: 'Service',
			resourceName: 'redis-cache',
			environment: 'production',
			details: 'Created new Redis service (1GB memory)',
			status: 'success',
			ipAddress: '192.168.1.112',
			userAgent: 'Mozilla/5.0 (X11; Linux x86_64)'
		},
		{
			id: 'act_008',
			timestamp: '2024-10-10T13:55:00Z',
			user: 'system',
			displayName: 'System',
			action: 'update',
			resource: 'Certificate',
			resourceName: 'app.mikrocloud.io',
			environment: 'production',
			details: 'Renewed SSL certificate (expires 2025-01-08)',
			status: 'success',
			ipAddress: 'internal',
			userAgent: 'mikrocloud-scheduler/1.0'
		},
		{
			id: 'act_009',
			timestamp: '2024-10-10T13:50:33Z',
			user: 'jane.smith',
			displayName: 'Jane Smith',
			action: 'deploy',
			resource: 'Application',
			resourceName: 'api-service',
			environment: 'staging',
			details: 'Deployed version v3.1.0 to staging',
			status: 'failed',
			ipAddress: '192.168.1.105',
			userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)',
			errorMessage: 'Health check failed: Connection refused on port 8080'
		},
		{
			id: 'act_010',
			timestamp: '2024-10-10T13:45:18Z',
			user: 'john.doe',
			displayName: 'John Doe',
			action: 'update',
			resource: 'Environment Variable',
			resourceName: 'DATABASE_URL',
			environment: 'production',
			details: 'Updated environment variable for mikrocloud-app',
			status: 'success',
			ipAddress: '192.168.1.100',
			userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)'
		},
		{
			id: 'act_011',
			timestamp: '2024-10-10T13:40:05Z',
			user: 'bob.johnson',
			displayName: 'Bob Johnson',
			action: 'delete',
			resource: 'Container',
			resourceName: 'old-worker-instance',
			environment: 'staging',
			details: 'Removed stopped container (unused for 30 days)',
			status: 'success',
			ipAddress: '192.168.1.112',
			userAgent: 'Mozilla/5.0 (X11; Linux x86_64)'
		},
		{
			id: 'act_012',
			timestamp: '2024-10-10T13:35:42Z',
			user: 'system',
			displayName: 'System',
			action: 'create',
			resource: 'Webhook',
			resourceName: 'github-deploy-hook',
			environment: null,
			details: 'Registered webhook for automatic deployments',
			status: 'success',
			ipAddress: 'internal',
			userAgent: 'mikrocloud-api/2.4.1'
		}
	];

	const filteredActivities = $derived(
		activities.filter((activity) => {
			const matchesSearch =
				activity.displayName.toLowerCase().includes(searchQuery.toLowerCase()) ||
				activity.resource.toLowerCase().includes(searchQuery.toLowerCase()) ||
				(activity.resourceName?.toLowerCase().includes(searchQuery.toLowerCase()) ?? false) ||
				activity.details.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesAction = actionFilter === 'all' || activity.action === actionFilter;
			const matchesUser = userFilter === 'all' || activity.user === userFilter;
			return matchesSearch && matchesAction && matchesUser;
		})
	);

	const stats = $derived({
		total: activities.length,
		today: activities.filter((a) => {
			const activityDate = new Date(a.timestamp);
			const today = new Date();
			return activityDate.toDateString() === today.toDateString();
		}).length,
		users: new Set(activities.map((a) => a.user)).size,
		failed: activities.filter((a) => a.status === 'failed').length
	});

	function getActionIcon(action: string) {
		switch (action) {
			case 'create':
				return FileEdit;
			case 'update':
				return Settings;
			case 'delete':
				return Trash2;
			case 'deploy':
				return GitBranch;
			case 'login':
				return Unlock;
			case 'permission':
				return Shield;
			default:
				return Database;
		}
	}

	function getActionColor(action: string) {
		switch (action) {
			case 'create':
				return 'text-green-500';
			case 'update':
				return 'text-blue-500';
			case 'delete':
				return 'text-red-500';
			case 'deploy':
				return 'text-purple-500';
			case 'login':
				return 'text-yellow-500';
			case 'permission':
				return 'text-orange-500';
			default:
				return 'text-gray-500';
		}
	}

	function getActionBgColor(action: string) {
		switch (action) {
			case 'create':
				return 'bg-green-500/10 border-green-500/20';
			case 'update':
				return 'bg-blue-500/10 border-blue-500/20';
			case 'delete':
				return 'bg-red-500/10 border-red-500/20';
			case 'deploy':
				return 'bg-purple-500/10 border-purple-500/20';
			case 'login':
				return 'bg-yellow-500/10 border-yellow-500/20';
			case 'permission':
				return 'bg-orange-500/10 border-orange-500/20';
			default:
				return 'bg-gray-500/10 border-gray-500/20';
		}
	}

	function formatTime(timestamp: string) {
		const date = new Date(timestamp);
		return date.toLocaleTimeString();
	}

	function formatDate(timestamp: string) {
		const date = new Date(timestamp);
		return date.toLocaleDateString();
	}

	function getRelativeTime(timestamp: string) {
		const date = new Date(timestamp);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);

		if (diffMins < 1) return 'Just now';
		if (diffMins < 60) return `${diffMins}m ago`;

		const diffHours = Math.floor(diffMins / 60);
		if (diffHours < 24) return `${diffHours}h ago`;

		const diffDays = Math.floor(diffHours / 24);
		return `${diffDays}d ago`;
	}

	async function handleRefresh() {
		isRefreshing = true;
		await new Promise((resolve) => setTimeout(resolve, 1000));
		isRefreshing = false;
	}

	function exportActivities() {
		const dataStr = JSON.stringify(filteredActivities, null, 2);
		const blob = new Blob([dataStr], { type: 'application/json' });
		const url = URL.createObjectURL(blob);
		const link = document.createElement('a');
		link.href = url;
		link.download = `activities-${new Date().toISOString()}.json`;
		link.click();
		URL.revokeObjectURL(url);
	}
</script>

<div class="flex-1 overflow-y-auto">
	<div class="border-b border-border bg-background sticky top-0 z-10">
		<div class="px-8 py-6">
			<div class="flex items-center justify-between mb-4">
				<div>
					<h1 class="text-3xl font-bold">Activities</h1>
					<p class="text-muted-foreground mt-1">Audit log of all user and system actions</p>
				</div>
				<div class="flex items-center gap-2">
					<Button variant="outline" size="sm" onclick={exportActivities}>Export</Button>
					<Button variant="outline" size="icon" onclick={handleRefresh} disabled={isRefreshing}>
						<RefreshCw class={`h-4 w-4 ${isRefreshing ? 'animate-spin' : ''} `} />
					</Button>
				</div>
			</div>

			<div class="grid grid-cols-4 gap-4 mb-4">
				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Total Activities</CardDescription>
						<CardTitle class="text-2xl">{stats.total}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Today</CardDescription>
						<CardTitle class="text-2xl text-blue-500">{stats.today}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Active Users</CardDescription>
						<CardTitle class="text-2xl text-green-500">{stats.users}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Failed Actions</CardDescription>
						<CardTitle class="text-2xl text-red-500">{stats.failed}</CardTitle>
					</CardHeader>
				</Card>
			</div>

			<div class="grid grid-cols-3 gap-3">
				<div class="relative">
					<Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
					<Input
						type="text"
						placeholder="Search activities..."
						class="pl-10"
						bind:value={searchQuery}
					/>
				</div>

				<Select.Root type="single" bind:value={actionFilter}>
					<Select.Trigger>
						{actionFilter || 'Filter by action'}
					</Select.Trigger>
					<Select.Content>
						{#each actionFilters as filter}
							<Select.Item value={filter.value}>{filter.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>

				<Select.Root type="single" bind:value={userFilter}>
					<Select.Trigger>
						{userFilter || 'Filter by user'}
					</Select.Trigger>
					<Select.Content>
						{#each userFilters as filter}
							<Select.Item value={filter.value}>{filter.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		</div>
	</div>

	<div class="p-8 space-y-3">
		<div class="text-sm text-muted-foreground mb-2">
			Showing {filteredActivities.length} of {activities.length} activities
		</div>

		{#each filteredActivities as activity (activity.id)}
			{@const Icon = getActionIcon(activity.action)}
			<Card class="hover:bg-muted/50 transition-colors">
				<CardContent class="p-4">
					<div class="flex items-start gap-4">
						<div class="p-2 border rounded-lg {getActionBgColor(activity.action)}">
							<Icon class="h-5 w-5 {getActionColor(activity.action)}" />
						</div>

						<div class="flex-1 min-w-0">
							<div class="flex items-start justify-between gap-4 mb-2">
								<div class="flex-1">
									<div class="flex items-center gap-2 mb-1">
										<span class="font-semibold capitalize">{activity.action}</span>
										<span class="text-muted-foreground">·</span>
										<span class="text-sm text-muted-foreground">{activity.resource}</span>
										{#if activity.resourceName}
											<span class="text-muted-foreground">·</span>
											<span class="text-sm font-mono">{activity.resourceName}</span>
										{/if}
										{#if activity.environment}
											<Badge variant="outline" class="text-xs">{activity.environment}</Badge>
										{/if}
									</div>
									<p class="text-sm text-muted-foreground mb-2">{activity.details}</p>

									{#if activity.status === 'failed' && activity.errorMessage}
										<div
											class="p-2 bg-red-500/10 border border-red-500/20 rounded text-sm text-red-500"
										>
											<span class="font-semibold">Error:</span>
											{activity.errorMessage}
										</div>
									{/if}
								</div>

								<div class="flex items-center gap-2">
									<Badge variant={activity.status === 'success' ? 'default' : 'destructive'}>
										{activity.status}
									</Badge>
								</div>
							</div>

							<div class="flex items-center gap-4 text-xs text-muted-foreground">
								<span class="flex items-center gap-1">
									<User class="h-3 w-3" />
									{activity.displayName}
								</span>
								<span class="flex items-center gap-1">
									<Calendar class="h-3 w-3" />
									{formatDate(activity.timestamp)} at {formatTime(activity.timestamp)}
								</span>
								<span>
									{getRelativeTime(activity.timestamp)}
								</span>
								<span class="font-mono">
									{activity.ipAddress}
								</span>
							</div>
						</div>
					</div>
				</CardContent>
			</Card>
		{/each}

		{#if filteredActivities.length === 0}
			<Card>
				<CardContent class="py-12">
					<div class="text-center text-muted-foreground">
						No activities found matching your criteria
					</div>
				</CardContent>
			</Card>
		{/if}
	</div>
</div>
