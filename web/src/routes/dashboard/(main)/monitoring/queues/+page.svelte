<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import { Search, RefreshCw, Clock, CheckCircle, AlertCircle, XCircle, RotateCw, Pause, Play } from 'lucide-svelte';

	let searchQuery = $state('');
	let statusFilter = $state('all');
	let queueFilter = $state('all');
	let isRefreshing = $state(false);

	const statusFilters = [
		{ value: 'all', label: 'All Statuses' },
		{ value: 'pending', label: 'Pending' },
		{ value: 'processing', label: 'Processing' },
		{ value: 'completed', label: 'Completed' },
		{ value: 'failed', label: 'Failed' },
		{ value: 'retry', label: 'Retry' }
	];

	const queueFilters = [
		{ value: 'all', label: 'All Queues' },
		{ value: 'deployments', label: 'Deployments' },
		{ value: 'builds', label: 'Builds' },
		{ value: 'backups', label: 'Backups' },
		{ value: 'webhooks', label: 'Webhooks' },
		{ value: 'emails', label: 'Emails' }
	];

	const jobs = [
		{
			id: 'job_001',
			queue: 'deployments',
			type: 'deploy_application',
			status: 'processing',
			priority: 1,
			payload: { appId: 'mikrocloud-app', version: 'v2.4.1', environment: 'production' },
			createdAt: '2024-10-10T14:35:00Z',
			startedAt: '2024-10-10T14:35:02Z',
			completedAt: null,
			retries: 0,
			maxRetries: 3,
			progress: 45,
			worker: 'worker-01'
		},
		{
			id: 'job_002',
			queue: 'builds',
			type: 'build_image',
			status: 'completed',
			priority: 2,
			payload: { appId: 'web-dashboard', branch: 'main', commit: 'a1b2c3d' },
			createdAt: '2024-10-10T14:30:00Z',
			startedAt: '2024-10-10T14:30:05Z',
			completedAt: '2024-10-10T14:33:45Z',
			retries: 0,
			maxRetries: 3,
			progress: 100,
			worker: 'worker-02',
			duration: 220
		},
		{
			id: 'job_003',
			queue: 'backups',
			type: 'database_backup',
			status: 'failed',
			priority: 1,
			payload: { dbId: 'postgres-main', type: 'full' },
			createdAt: '2024-10-10T14:25:00Z',
			startedAt: '2024-10-10T14:25:03Z',
			completedAt: '2024-10-10T14:27:15Z',
			retries: 2,
			maxRetries: 3,
			progress: 78,
			worker: 'worker-03',
			duration: 132,
			error: 'Connection timeout: Failed to reach database endpoint'
		},
		{
			id: 'job_004',
			queue: 'webhooks',
			type: 'send_webhook',
			status: 'pending',
			priority: 3,
			payload: { event: 'deployment.success', url: 'https://hooks.example.com/deploy' },
			createdAt: '2024-10-10T14:34:00Z',
			startedAt: null,
			completedAt: null,
			retries: 0,
			maxRetries: 5,
			progress: 0,
			worker: null
		},
		{
			id: 'job_005',
			queue: 'emails',
			type: 'send_notification',
			status: 'completed',
			priority: 4,
			payload: { to: 'user@example.com', template: 'deployment_success' },
			createdAt: '2024-10-10T14:20:00Z',
			startedAt: '2024-10-10T14:20:01Z',
			completedAt: '2024-10-10T14:20:05Z',
			retries: 0,
			maxRetries: 3,
			progress: 100,
			worker: 'worker-01',
			duration: 4
		},
		{
			id: 'job_006',
			queue: 'deployments',
			type: 'rollback_deployment',
			status: 'retry',
			priority: 1,
			payload: { appId: 'api-service', version: 'v3.0.9' },
			createdAt: '2024-10-10T14:15:00Z',
			startedAt: '2024-10-10T14:15:05Z',
			completedAt: null,
			retries: 1,
			maxRetries: 3,
			progress: 0,
			worker: null,
			error: 'Service health check failed, retrying...'
		},
		{
			id: 'job_007',
			queue: 'builds',
			type: 'cleanup_old_images',
			status: 'completed',
			priority: 5,
			payload: { olderThan: '30d', namespace: 'production' },
			createdAt: '2024-10-10T14:10:00Z',
			startedAt: '2024-10-10T14:10:02Z',
			completedAt: '2024-10-10T14:11:30Z',
			retries: 0,
			maxRetries: 3,
			progress: 100,
			worker: 'worker-02',
			duration: 88
		},
		{
			id: 'job_008',
			queue: 'backups',
			type: 'volume_snapshot',
			status: 'processing',
			priority: 2,
			payload: { volumeId: 'vol_data_01', compression: 'gzip' },
			createdAt: '2024-10-10T14:32:00Z',
			startedAt: '2024-10-10T14:32:10Z',
			completedAt: null,
			retries: 0,
			maxRetries: 3,
			progress: 62,
			worker: 'worker-03'
		}
	];

	const filteredJobs = $derived(
		jobs.filter((job) => {
			const matchesSearch =
				job.type.toLowerCase().includes(searchQuery.toLowerCase()) ||
				job.queue.toLowerCase().includes(searchQuery.toLowerCase()) ||
				job.id.toLowerCase().includes(searchQuery.toLowerCase()) ||
				JSON.stringify(job.payload).toLowerCase().includes(searchQuery.toLowerCase());
			const matchesStatus = statusFilter === 'all' || job.status === statusFilter;
			const matchesQueue = queueFilter === 'all' || job.queue === queueFilter;
			return matchesSearch && matchesStatus && matchesQueue;
		})
	);

	const stats = $derived({
		total: jobs.length,
		pending: jobs.filter((j) => j.status === 'pending').length,
		processing: jobs.filter((j) => j.status === 'processing').length,
		completed: jobs.filter((j) => j.status === 'completed').length,
		failed: jobs.filter((j) => j.status === 'failed').length,
		retry: jobs.filter((j) => j.status === 'retry').length,
		avgDuration: Math.round(
			jobs.filter((j) => j.duration).reduce((sum, j) => sum + (j.duration || 0), 0) /
				jobs.filter((j) => j.duration).length || 0
		)
	});

	function getStatusBadgeVariant(status: string): 'default' | 'secondary' | 'destructive' {
		switch (status) {
			case 'completed':
				return 'default';
			case 'failed':
				return 'destructive';
			case 'processing':
			case 'pending':
			case 'retry':
				return 'secondary';
			default:
				return 'secondary';
		}
	}

	function getStatusIcon(status: string) {
		switch (status) {
			case 'completed':
				return CheckCircle;
			case 'failed':
				return XCircle;
			case 'processing':
				return Play;
			case 'retry':
				return RotateCw;
			case 'pending':
				return Pause;
			default:
				return Clock;
		}
	}

	function getStatusColor(status: string) {
		switch (status) {
			case 'completed':
				return 'text-green-500';
			case 'failed':
				return 'text-red-500';
			case 'processing':
				return 'text-blue-500';
			case 'retry':
				return 'text-orange-500';
			case 'pending':
				return 'text-yellow-500';
			default:
				return 'text-gray-500';
		}
	}

	function getPriorityBadge(priority: number) {
		if (priority <= 2) return { label: 'High', variant: 'destructive' as const };
		if (priority <= 4) return { label: 'Medium', variant: 'secondary' as const };
		return { label: 'Low', variant: 'default' as const };
	}

	function formatDuration(seconds: number | null) {
		if (seconds === null) return '-';
		if (seconds < 60) return `${seconds}s`;
		const minutes = Math.floor(seconds / 60);
		const secs = seconds % 60;
		return `${minutes}m ${secs}s`;
	}

	function formatTime(timestamp: string | null) {
		if (!timestamp) return '-';
		const date = new Date(timestamp);
		return date.toLocaleTimeString();
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

	function retryJob(jobId: string) {
		console.log('Retrying job:', jobId);
	}

	function cancelJob(jobId: string) {
		console.log('Cancelling job:', jobId);
	}
</script>

<div class="flex-1 overflow-y-auto">
	<div class="border-b border-border bg-background sticky top-0 z-10">
		<div class="px-8 py-6">
			<div class="flex items-center justify-between mb-4">
				<div>
					<h1 class="text-3xl font-bold">Queues</h1>
					<p class="text-muted-foreground mt-1">Monitor job queues, processing status, and retry management</p>
				</div>
			<Button variant="outline" size="icon" onclick={handleRefresh} disabled={isRefreshing}>
				<RefreshCw class={`h-4 w-4 ${isRefreshing ? 'animate-spin' : ''}`} />
			</Button>
			</div>

			<div class="grid grid-cols-5 gap-4 mb-4">
				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Total Jobs</CardDescription>
						<CardTitle class="text-2xl">{stats.total}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Pending</CardDescription>
						<CardTitle class="text-2xl text-yellow-500">{stats.pending}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Processing</CardDescription>
						<CardTitle class="text-2xl text-blue-500">{stats.processing}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Failed</CardDescription>
						<CardTitle class="text-2xl text-red-500">{stats.failed}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Avg Duration</CardDescription>
						<CardTitle class="text-2xl">{formatDuration(stats.avgDuration)}</CardTitle>
					</CardHeader>
				</Card>
			</div>

			<div class="grid grid-cols-3 gap-3">
				<div class="relative">
					<Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
					<Input
						type="text"
						placeholder="Search jobs..."
						class="pl-10"
						bind:value={searchQuery}
					/>
				</div>

				<Select.Root
					selected={{ value: statusFilter, label: statusFilter }}
					onSelectedChange={(v) => v && (statusFilter = v.value)}
				>
					<Select.Trigger>
						<Select.Value placeholder="Filter by status" />
					</Select.Trigger>
					<Select.Content>
						{#each statusFilters as filter}
							<Select.Item value={filter.value}>{filter.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>

				<Select.Root
					selected={{ value: queueFilter, label: queueFilter }}
					onSelectedChange={(v) => v && (queueFilter = v.value)}
				>
					<Select.Trigger>
						<Select.Value placeholder="Filter by queue" />
					</Select.Trigger>
					<Select.Content>
						{#each queueFilters as filter}
							<Select.Item value={filter.value}>{filter.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		</div>
	</div>

	<div class="p-8 space-y-4">
		<div class="text-sm text-muted-foreground mb-2">
			Showing {filteredJobs.length} of {jobs.length} jobs
		</div>

	{#each filteredJobs as job (job.id)}
		{@const priorityBadge = getPriorityBadge(job.priority)}
		{@const Icon = getStatusIcon(job.status)}
		<Card>
			<CardHeader>
				<div class="flex items-start justify-between">
					<div class="flex-1">
						<div class="flex items-center gap-3 mb-2">
							<CardTitle class="text-lg font-mono">{job.id}</CardTitle>
							<Badge variant="outline">{job.queue}</Badge>
							<Badge variant={getStatusBadgeVariant(job.status)}>
								{job.status}
							</Badge>
							<Badge variant={priorityBadge.variant}>
								{priorityBadge.label}
							</Badge>
						</div>
							<CardDescription class="flex items-center gap-4 text-sm">
								<span class="font-medium">{job.type}</span>
								<span class="text-muted-foreground">
									Created {getRelativeTime(job.createdAt)}
								</span>
								{#if job.worker}
									<span class="text-muted-foreground">
										Worker: {job.worker}
									</span>
								{/if}
							</CardDescription>
						</div>

					<div class="flex items-center gap-2">
						<Icon class="h-5 w-5 {getStatusColor(job.status)}" />
					</div>
					</div>
				</CardHeader>
				<CardContent>
					<div class="space-y-4">
						{#if job.status === 'processing' || job.status === 'retry'}
							<div>
								<div class="flex items-center justify-between text-sm mb-2">
									<span class="text-muted-foreground">Progress</span>
									<span class="font-medium">{job.progress}%</span>
								</div>
								<div class="w-full bg-muted rounded-full h-2">
									<div
										class="bg-blue-500 h-2 rounded-full transition-all duration-300"
										style="width: {job.progress}%"
									></div>
								</div>
							</div>
						{/if}

						<div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
							<div>
								<div class="text-xs text-muted-foreground mb-1">Created</div>
								<div class="text-sm font-medium">{formatTime(job.createdAt)}</div>
							</div>

							<div>
								<div class="text-xs text-muted-foreground mb-1">Started</div>
								<div class="text-sm font-medium">{formatTime(job.startedAt)}</div>
							</div>

							<div>
								<div class="text-xs text-muted-foreground mb-1">Completed</div>
								<div class="text-sm font-medium">{formatTime(job.completedAt)}</div>
							</div>

							<div>
								<div class="text-xs text-muted-foreground mb-1">Duration</div>
								<div class="text-sm font-medium flex items-center gap-1">
									<Clock class="h-3 w-3" />
									{formatDuration(job.duration ?? null)}
								</div>
							</div>
						</div>

						<div>
							<div class="text-xs text-muted-foreground mb-2">Retry Status</div>
							<div class="flex items-center gap-4 text-sm">
								<span>
									Attempts: <span class="font-medium">{job.retries}/{job.maxRetries}</span>
								</span>
								<div class="flex gap-1">
									{#each Array(job.maxRetries) as _, i}
										<div
											class="w-2 h-2 rounded-full {i < job.retries ? 'bg-orange-500' : 'bg-muted'}"
										></div>
									{/each}
								</div>
							</div>
						</div>

						{#if job.error}
							<div class="p-3 bg-red-500/10 border border-red-500/20 rounded-lg">
								<div class="flex items-start gap-2">
									<AlertCircle class="h-4 w-4 text-red-500 mt-0.5" />
									<div class="flex-1">
										<div class="text-sm font-semibold text-red-500 mb-1">Error</div>
										<div class="text-sm text-red-500/90">{job.error}</div>
									</div>
								</div>
							</div>
						{/if}

						<details class="group">
							<summary class="cursor-pointer text-sm text-muted-foreground hover:text-foreground flex items-center gap-2">
								<span>Payload</span>
								<svg
									class="w-4 h-4 transition-transform group-open:rotate-180"
									fill="none"
									stroke="currentColor"
									viewBox="0 0 24 24"
								>
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
								</svg>
							</summary>
							<div class="mt-2 p-3 bg-muted rounded-lg">
								<pre class="text-xs font-mono overflow-x-auto">{JSON.stringify(job.payload, null, 2)}</pre>
							</div>
						</details>

						{#if job.status === 'failed' || job.status === 'retry'}
							<div class="pt-4 border-t flex items-center justify-end gap-2">
								<Button variant="outline" size="sm" onclick={() => retryJob(job.id)}>
									<RotateCw class="h-3 w-3 mr-1" />
									Retry Now
								</Button>
								<Button variant="outline" size="sm" onclick={() => cancelJob(job.id)}>
									Cancel
								</Button>
							</div>
						{:else if job.status === 'pending' || job.status === 'processing'}
							<div class="pt-4 border-t flex items-center justify-end gap-2">
								<Button variant="outline" size="sm" onclick={() => cancelJob(job.id)}>
									Cancel Job
								</Button>
							</div>
						{/if}
					</div>
				</CardContent>
			</Card>
		{/each}

		{#if filteredJobs.length === 0}
			<Card>
				<CardContent class="py-12">
					<div class="text-center text-muted-foreground">
						No jobs found matching your criteria
					</div>
				</CardContent>
			</Card>
		{/if}
	</div>
</div>
