<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import { Search, RefreshCw, Download, Calendar, Filter, Terminal } from 'lucide-svelte';

	let searchQuery = $state('');
	let levelFilter = $state('all');
	let sourceFilter = $state('all');
	let timeRange = $state('last-hour');
	let isRefreshing = $state(false);
	let autoScroll = $state(true);

	const logLevels = [
		{ value: 'all', label: 'All Levels' },
		{ value: 'debug', label: 'Debug' },
		{ value: 'info', label: 'Info' },
		{ value: 'warn', label: 'Warning' },
		{ value: 'error', label: 'Error' },
		{ value: 'fatal', label: 'Fatal' }
	];

	const logSources = [
		{ value: 'all', label: 'All Sources' },
		{ value: 'api', label: 'API Server' },
		{ value: 'scheduler', label: 'Scheduler' },
		{ value: 'postgres', label: 'PostgreSQL' },
		{ value: 'app-web', label: 'App: web' },
		{ value: 'traefik', label: 'Traefik' }
	];

	const timeRanges = [
		{ value: 'last-5-min', label: 'Last 5 minutes' },
		{ value: 'last-15-min', label: 'Last 15 minutes' },
		{ value: 'last-hour', label: 'Last hour' },
		{ value: 'last-6-hours', label: 'Last 6 hours' },
		{ value: 'last-24-hours', label: 'Last 24 hours' },
		{ value: 'last-7-days', label: 'Last 7 days' }
	];

	const logs = [
		{
			id: 1,
			timestamp: '2024-10-10T14:32:15.123Z',
			level: 'info',
			source: 'api',
			message: 'Successfully created deployment for project mikrocloud-app',
			metadata: { deploymentId: 'dep_abc123', projectId: 'proj_xyz789' }
		},
		{
			id: 2,
			timestamp: '2024-10-10T14:32:14.856Z',
			level: 'debug',
			source: 'scheduler',
			message: 'Checking for pending jobs in queue',
			metadata: { queueSize: 12, workers: 3 }
		},
		{
			id: 3,
			timestamp: '2024-10-10T14:32:13.421Z',
			level: 'info',
			source: 'api',
			message: 'HTTP POST /api/v1/deployments - 201 Created (45ms)',
			metadata: { method: 'POST', path: '/api/v1/deployments', status: 201, duration: 45 }
		},
		{
			id: 4,
			timestamp: '2024-10-10T14:32:12.189Z',
			level: 'warn',
			source: 'postgres',
			message: 'Slow query detected: SELECT * FROM applications WHERE...',
			metadata: { duration: 1234, queryType: 'SELECT' }
		},
		{
			id: 5,
			timestamp: '2024-10-10T14:32:10.567Z',
			level: 'error',
			source: 'app-web',
			message: 'Failed to connect to database: connection timeout',
			metadata: { error: 'ETIMEDOUT', host: 'postgres-main:5432' }
		},
		{
			id: 6,
			timestamp: '2024-10-10T14:32:09.234Z',
			level: 'info',
			source: 'traefik',
			message: 'Request handled: GET /dashboard -> 200 OK (12ms)',
			metadata: { method: 'GET', path: '/dashboard', status: 200, duration: 12 }
		},
		{
			id: 7,
			timestamp: '2024-10-10T14:32:08.890Z',
			level: 'debug',
			source: 'api',
			message: 'Cache hit for user session: user_456',
			metadata: { userId: 'user_456', cacheKey: 'session:user_456' }
		},
		{
			id: 8,
			timestamp: '2024-10-10T14:32:07.456Z',
			level: 'info',
			source: 'scheduler',
			message: 'Completed job: backup-database',
			metadata: { jobId: 'job_backup_123', duration: 5420 }
		},
		{
			id: 9,
			timestamp: '2024-10-10T14:32:06.123Z',
			level: 'error',
			source: 'api',
			message: 'Authentication failed: invalid token',
			metadata: { token: '***', ip: '192.168.1.45' }
		},
		{
			id: 10,
			timestamp: '2024-10-10T14:32:05.789Z',
			level: 'fatal',
			source: 'app-web',
			message: 'Application crashed due to unhandled exception: TypeError',
			metadata: { error: 'TypeError: Cannot read property of undefined', stack: '...' }
		},
		{
			id: 11,
			timestamp: '2024-10-10T14:32:04.456Z',
			level: 'info',
			source: 'api',
			message: 'Server started on port 8080',
			metadata: { port: 8080, environment: 'production' }
		},
		{
			id: 12,
			timestamp: '2024-10-10T14:32:03.123Z',
			level: 'debug',
			source: 'scheduler',
			message: 'Initialized worker pool with 5 workers',
			metadata: { workerCount: 5 }
		}
	];

	const filteredLogs = $derived(
		logs.filter((log) => {
			const matchesSearch =
				log.message.toLowerCase().includes(searchQuery.toLowerCase()) ||
				log.source.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesLevel = levelFilter === 'all' || log.level === levelFilter;
			const matchesSource = sourceFilter === 'all' || log.source === sourceFilter;
			return matchesSearch && matchesLevel && matchesSource;
		})
	);

	function getLevelColor(level: string) {
		switch (level) {
			case 'debug':
				return 'text-gray-500';
			case 'info':
				return 'text-blue-500';
			case 'warn':
				return 'text-yellow-500';
			case 'error':
				return 'text-red-500';
			case 'fatal':
				return 'text-red-700 font-bold';
			default:
				return 'text-gray-500';
		}
	}

	function getLevelBadgeVariant(level: string): 'default' | 'secondary' | 'destructive' {
		switch (level) {
			case 'error':
			case 'fatal':
				return 'destructive';
			case 'warn':
				return 'secondary';
			default:
				return 'default';
		}
	}

	function formatTime(timestamp: string) {
		const date = new Date(timestamp);
		return date.toLocaleTimeString() + '.' + date.getMilliseconds();
	}

	async function handleRefresh() {
		isRefreshing = true;
		await new Promise((resolve) => setTimeout(resolve, 1000));
		isRefreshing = false;
	}

	function handleExport() {
		const dataStr = JSON.stringify(filteredLogs, null, 2);
		const blob = new Blob([dataStr], { type: 'application/json' });
		const url = URL.createObjectURL(blob);
		const link = document.createElement('a');
		link.href = url;
		link.download = `logs-${new Date().toISOString()}.json`;
		link.click();
		URL.revokeObjectURL(url);
	}
</script>

<div class="flex-1 overflow-y-auto">
	<div class="border-b border-border bg-background sticky top-0 z-10">
		<div class="px-8 py-6">
			<div class="flex items-center justify-between mb-4">
				<div>
					<h1 class="text-3xl font-bold">Logs</h1>
					<p class="text-muted-foreground mt-1">
						Real-time log streaming and search across all services
					</p>
				</div>
				<div class="flex items-center gap-2">
					<Button variant="outline" size="sm" onclick={handleExport}>
						<Download class="h-4 w-4 mr-2" />
						Export
					</Button>
				<Button variant="outline" size="icon" onclick={handleRefresh} disabled={isRefreshing}>
					<RefreshCw class={`h-4 w-4 ${isRefreshing ? 'animate-spin' : ''}`} />
				</Button>
				</div>
			</div>

			<div class="grid grid-cols-4 gap-3">
				<div class="relative col-span-2">
					<Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
					<Input
						type="text"
						placeholder="Search logs..."
						class="pl-10"
						bind:value={searchQuery}
					/>
				</div>

				<Select.Root
					selected={{ value: levelFilter, label: levelFilter }}
					onSelectedChange={(v) => v && (levelFilter = v.value)}
				>
					<Select.Trigger>
						<Filter class="h-4 w-4 mr-2" />
						<Select.Value placeholder="Level" />
					</Select.Trigger>
					<Select.Content>
						{#each logLevels as level}
							<Select.Item value={level.value}>{level.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>

				<Select.Root
					selected={{ value: sourceFilter, label: sourceFilter }}
					onSelectedChange={(v) => v && (sourceFilter = v.value)}
				>
					<Select.Trigger>
						<Terminal class="h-4 w-4 mr-2" />
						<Select.Value placeholder="Source" />
					</Select.Trigger>
					<Select.Content>
						{#each logSources as source}
							<Select.Item value={source.value}>{source.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>

			<div class="flex items-center gap-3 mt-3">
				<Select.Root
					selected={{ value: timeRange, label: timeRange }}
					onSelectedChange={(v) => v && (timeRange = v.value)}
				>
					<Select.Trigger class="w-[200px]">
						<Calendar class="h-4 w-4 mr-2" />
						<Select.Value placeholder="Time range" />
					</Select.Trigger>
					<Select.Content>
						{#each timeRanges as range}
							<Select.Item value={range.value}>{range.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>

				<div class="flex items-center gap-2 ml-auto">
					<label class="flex items-center gap-2 text-sm cursor-pointer">
						<input type="checkbox" bind:checked={autoScroll} class="rounded" />
						Auto-scroll
					</label>
					<Badge variant="secondary">{filteredLogs.length} logs</Badge>
				</div>
			</div>
		</div>
	</div>

	<div class="p-8">
		<Card>
			<CardHeader>
				<div class="flex items-center justify-between">
					<div>
						<CardTitle>Log Stream</CardTitle>
						<CardDescription>
							Showing logs from {timeRanges.find((r) => r.value === timeRange)?.label ||
								'selected range'}
						</CardDescription>
					</div>
				</div>
			</CardHeader>
			<CardContent class="p-0">
				<div class="divide-y max-h-[calc(100vh-400px)] overflow-y-auto font-mono text-sm">
					{#each filteredLogs as log (log.id)}
						<div class="px-6 py-3 hover:bg-accent/50 transition-colors">
							<div class="flex items-start gap-4">
								<div class="text-xs text-muted-foreground whitespace-nowrap pt-0.5 min-w-[140px]">
									{formatTime(log.timestamp)}
								</div>
								<div class="min-w-[80px]">
									<Badge variant={getLevelBadgeVariant(log.level)} class="text-xs uppercase">
										{log.level}
									</Badge>
								</div>
								<div class="min-w-[120px]">
									<Badge variant="outline" class="text-xs font-mono">
										{log.source}
									</Badge>
								</div>
								<div class="flex-1 {getLevelColor(log.level)} break-all">
									{log.message}
								</div>
							</div>
							{#if Object.keys(log.metadata).length > 0}
								<div class="mt-2 ml-[344px] text-xs text-muted-foreground">
									<details class="cursor-pointer">
										<summary class="hover:text-foreground">Show metadata</summary>
										<pre
											class="mt-2 p-3 bg-muted rounded-md overflow-x-auto">{JSON.stringify(log.metadata, null, 2)}</pre>
									</details>
								</div>
							{/if}
						</div>
					{/each}

					{#if filteredLogs.length === 0}
						<div class="py-12 text-center text-muted-foreground">
							No logs found matching your criteria
						</div>
					{/if}
				</div>
			</CardContent>
		</Card>

		<Card class="mt-4">
			<CardHeader>
				<CardTitle class="text-base">Log Retention Policy</CardTitle>
			</CardHeader>
			<CardContent>
				<div class="grid gap-4 md:grid-cols-3">
					<div>
						<div class="text-sm font-medium mb-1">Current Plan</div>
						<div class="text-2xl font-bold">7 Days</div>
						<div class="text-xs text-muted-foreground mt-1">Production environment</div>
					</div>
					<div>
						<div class="text-sm font-medium mb-1">Storage Used</div>
						<div class="text-2xl font-bold">2.4 GB</div>
						<div class="text-xs text-muted-foreground mt-1">of 10 GB limit</div>
						<div class="h-2 bg-secondary rounded-full overflow-hidden mt-2">
							<div class="h-full bg-primary" style="width: 24%"></div>
						</div>
					</div>
					<div>
						<div class="text-sm font-medium mb-1">Total Logs</div>
						<div class="text-2xl font-bold">1.2M</div>
						<div class="text-xs text-muted-foreground mt-1">in the last 7 days</div>
					</div>
				</div>
			</CardContent>
		</Card>
	</div>
</div>
