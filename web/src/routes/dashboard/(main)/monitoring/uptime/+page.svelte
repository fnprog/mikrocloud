<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import { Search, RefreshCw, Check, X, Clock, TrendingUp, Activity, AlertTriangle, Globe, ChevronDown, ChevronUp } from 'lucide-svelte';

	let searchQuery = $state('');
	let statusFilter = $state('all');
	let typeFilter = $state('all');
	let isRefreshing = $state(false);
	let expandedEndpoint = $state<string | null>(null);

	const statusFilters = [
		{ value: 'all', label: 'All Status' },
		{ value: 'up', label: 'Up' },
		{ value: 'down', label: 'Down' },
		{ value: 'degraded', label: 'Degraded' }
	];

	const typeFilters = [
		{ value: 'all', label: 'All Types' },
		{ value: 'http', label: 'HTTP' },
		{ value: 'https', label: 'HTTPS' },
		{ value: 'tcp', label: 'TCP' },
		{ value: 'ping', label: 'PING' }
	];

	const endpoints = [
		{
			id: 'ep_001',
			name: 'Main Website',
			url: 'https://mikrocloud.io',
			type: 'https',
			status: 'up',
			uptime: 99.98,
			responseTime: 142,
			lastChecked: '2024-10-10T14:35:00Z',
			checkInterval: 60,
			incidents: [
				{ timestamp: '2024-10-08T03:15:00Z', duration: 180, reason: 'Database connection timeout' }
			]
		},
		{
			id: 'ep_002',
			name: 'API Gateway',
			url: 'https://api.mikrocloud.io',
			type: 'https',
			status: 'up',
			uptime: 99.95,
			responseTime: 89,
			lastChecked: '2024-10-10T14:35:10Z',
			checkInterval: 30,
			incidents: [
				{ timestamp: '2024-10-09T14:30:00Z', duration: 420, reason: 'High load, rate limiting triggered' },
				{ timestamp: '2024-10-07T22:00:00Z', duration: 120, reason: 'Deployment rollout' }
			]
		},
		{
			id: 'ep_003',
			name: 'Database Cluster',
			url: 'postgres.internal:5432',
			type: 'tcp',
			status: 'up',
			uptime: 99.99,
			responseTime: 12,
			lastChecked: '2024-10-10T14:35:15Z',
			checkInterval: 60,
			incidents: []
		},
		{
			id: 'ep_004',
			name: 'Redis Cache',
			url: 'redis.internal:6379',
			type: 'tcp',
			status: 'up',
			uptime: 100.0,
			responseTime: 5,
			lastChecked: '2024-10-10T14:35:20Z',
			checkInterval: 30,
			incidents: []
		},
		{
			id: 'ep_005',
			name: 'Staging Environment',
			url: 'https://staging.mikrocloud.io',
			type: 'https',
			status: 'degraded',
			uptime: 97.5,
			responseTime: 2340,
			lastChecked: '2024-10-10T14:35:25Z',
			checkInterval: 120,
			incidents: [
				{ timestamp: '2024-10-10T12:00:00Z', duration: 3600, reason: 'Ongoing: Memory pressure, slow responses' },
				{ timestamp: '2024-10-09T08:15:00Z', duration: 900, reason: 'Container restart' },
				{ timestamp: '2024-10-08T16:45:00Z', duration: 1200, reason: 'Disk I/O bottleneck' }
			]
		},
		{
			id: 'ep_006',
			name: 'CDN Edge Node',
			url: 'https://cdn.mikrocloud.io',
			type: 'https',
			status: 'up',
			uptime: 99.92,
			responseTime: 28,
			lastChecked: '2024-10-10T14:35:30Z',
			checkInterval: 60,
			incidents: [
				{ timestamp: '2024-10-06T10:30:00Z', duration: 600, reason: 'CDN provider maintenance' }
			]
		},
		{
			id: 'ep_007',
			name: 'Webhook Service',
			url: 'https://webhooks.mikrocloud.io',
			type: 'https',
			status: 'up',
			uptime: 99.87,
			responseTime: 156,
			lastChecked: '2024-10-10T14:35:35Z',
			checkInterval: 120,
			incidents: [
				{ timestamp: '2024-10-09T19:20:00Z', duration: 240, reason: 'Queue overflow' },
				{ timestamp: '2024-10-05T04:00:00Z', duration: 480, reason: 'Network connectivity issues' }
			]
		},
		{
			id: 'ep_008',
			name: 'Legacy API',
			url: 'http://legacy-api.internal:8080',
			type: 'http',
			status: 'down',
			uptime: 85.3,
			responseTime: 0,
			lastChecked: '2024-10-10T14:35:40Z',
			checkInterval: 300,
			incidents: [
				{ timestamp: '2024-10-10T10:00:00Z', duration: null, reason: 'Ongoing: Connection refused' },
				{ timestamp: '2024-10-09T16:30:00Z', duration: 1800, reason: 'Service crash, OOM error' },
				{ timestamp: '2024-10-08T12:15:00Z', duration: 3600, reason: 'Manual shutdown for maintenance' },
				{ timestamp: '2024-10-07T08:00:00Z', duration: 2400, reason: 'Database migration failure' }
			]
		},
		{
			id: 'ep_009',
			name: 'Metrics Collector',
			url: 'metrics.internal',
			type: 'ping',
			status: 'up',
			uptime: 99.96,
			responseTime: 3,
			lastChecked: '2024-10-10T14:35:45Z',
			checkInterval: 60,
			incidents: [
				{ timestamp: '2024-10-04T02:30:00Z', duration: 300, reason: 'Host reboot' }
			]
		},
		{
			id: 'ep_010',
			name: 'Build Server',
			url: 'https://builds.mikrocloud.io',
			type: 'https',
			status: 'up',
			uptime: 99.85,
			responseTime: 234,
			lastChecked: '2024-10-10T14:35:50Z',
			checkInterval: 120,
			incidents: [
				{ timestamp: '2024-10-08T18:00:00Z', duration: 900, reason: 'Build queue overflow' },
				{ timestamp: '2024-10-03T14:20:00Z', duration: 600, reason: 'Storage full, cleanup required' }
			]
		}
	];

	const filteredEndpoints = $derived(
		endpoints.filter((endpoint) => {
			const matchesSearch =
				endpoint.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				endpoint.url.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesStatus = statusFilter === 'all' || endpoint.status === statusFilter;
			const matchesType = typeFilter === 'all' || endpoint.type === typeFilter;
			return matchesSearch && matchesStatus && matchesType;
		})
	);

	const stats = $derived({
		total: endpoints.length,
		up: endpoints.filter((e) => e.status === 'up').length,
		down: endpoints.filter((e) => e.status === 'down').length,
		degraded: endpoints.filter((e) => e.status === 'degraded').length,
		avgUptime: (endpoints.reduce((sum, e) => sum + e.uptime, 0) / endpoints.length).toFixed(2),
		avgResponseTime: Math.round(
			endpoints.filter((e) => e.status !== 'down').reduce((sum, e) => sum + e.responseTime, 0) /
				endpoints.filter((e) => e.status !== 'down').length
		)
	});

	function getStatusColor(status: string) {
		switch (status) {
			case 'up':
				return 'text-green-500';
			case 'down':
				return 'text-red-500';
			case 'degraded':
				return 'text-yellow-500';
			default:
				return 'text-gray-500';
		}
	}

	function getStatusBgColor(status: string) {
		switch (status) {
			case 'up':
				return 'bg-green-500/10 border-green-500/20';
			case 'down':
				return 'bg-red-500/10 border-red-500/20';
			case 'degraded':
				return 'bg-yellow-500/10 border-yellow-500/20';
			default:
				return 'bg-gray-500/10 border-gray-500/20';
		}
	}

	function getStatusIcon(status: string) {
		switch (status) {
			case 'up':
				return Check;
			case 'down':
				return X;
			case 'degraded':
				return AlertTriangle;
			default:
				return Activity;
		}
	}

	function formatResponseTime(ms: number) {
		if (ms === 0) return 'N/A';
		if (ms < 1000) return `${ms}ms`;
		return `${(ms / 1000).toFixed(2)}s`;
	}

	function formatDuration(seconds: number | null) {
		if (seconds === null) return 'Ongoing';
		if (seconds < 60) return `${seconds}s`;
		if (seconds < 3600) return `${Math.floor(seconds / 60)}m`;
		return `${Math.floor(seconds / 3600)}h ${Math.floor((seconds % 3600) / 60)}m`;
	}

	function formatTimestamp(timestamp: string) {
		const date = new Date(timestamp);
		return date.toLocaleString();
	}

	function getRelativeTime(timestamp: string) {
		const date = new Date(timestamp);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffSecs = Math.floor(diffMs / 1000);
		
		if (diffSecs < 60) return `${diffSecs}s ago`;
		
		const diffMins = Math.floor(diffSecs / 60);
		if (diffMins < 60) return `${diffMins}m ago`;
		
		const diffHours = Math.floor(diffMins / 60);
		if (diffHours < 24) return `${diffHours}h ago`;
		
		const diffDays = Math.floor(diffHours / 24);
		return `${diffDays}d ago`;
	}

	function toggleExpanded(id: string) {
		expandedEndpoint = expandedEndpoint === id ? null : id;
	}

	async function handleRefresh() {
		isRefreshing = true;
		await new Promise((resolve) => setTimeout(resolve, 1000));
		isRefreshing = false;
	}

	function exportEndpoints() {
		const dataStr = JSON.stringify(filteredEndpoints, null, 2);
		const blob = new Blob([dataStr], { type: 'application/json' });
		const url = URL.createObjectURL(blob);
		const link = document.createElement('a');
		link.href = url;
		link.download = `uptime-${new Date().toISOString()}.json`;
		link.click();
		URL.revokeObjectURL(url);
	}
</script>

<div class="flex-1 overflow-y-auto">
	<div class="border-b border-border bg-background sticky top-0 z-10">
		<div class="px-8 py-6">
			<div class="flex items-center justify-between mb-4">
				<div>
					<h1 class="text-3xl font-bold">Uptime Monitoring</h1>
					<p class="text-muted-foreground mt-1">Monitor endpoint health and availability</p>
				</div>
				<div class="flex items-center gap-2">
					<Button variant="outline" size="sm" onclick={exportEndpoints}>
						Export
					</Button>
				<Button variant="outline" size="icon" onclick={handleRefresh} disabled={isRefreshing}>
					<RefreshCw class={`h-4 w-4 ${isRefreshing ? 'animate-spin' : ''}`} />
				</Button>
				</div>
			</div>

			<div class="grid grid-cols-6 gap-4 mb-4">
				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Total Endpoints</CardDescription>
						<CardTitle class="text-2xl">{stats.total}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Up</CardDescription>
						<CardTitle class="text-2xl text-green-500">{stats.up}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Down</CardDescription>
						<CardTitle class="text-2xl text-red-500">{stats.down}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Degraded</CardDescription>
						<CardTitle class="text-2xl text-yellow-500">{stats.degraded}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Avg Uptime</CardDescription>
						<CardTitle class="text-2xl text-blue-500">{stats.avgUptime}%</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Avg Response</CardDescription>
						<CardTitle class="text-2xl text-purple-500">{stats.avgResponseTime}ms</CardTitle>
					</CardHeader>
				</Card>
			</div>

			<div class="grid grid-cols-3 gap-3">
				<div class="relative">
					<Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
					<Input
						type="text"
						placeholder="Search endpoints..."
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
					selected={{ value: typeFilter, label: typeFilter }}
					onSelectedChange={(v) => v && (typeFilter = v.value)}
				>
					<Select.Trigger>
						<Select.Value placeholder="Filter by type" />
					</Select.Trigger>
					<Select.Content>
						{#each typeFilters as filter}
							<Select.Item value={filter.value}>{filter.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		</div>
	</div>

	<div class="p-8 space-y-3">
		<div class="text-sm text-muted-foreground mb-2">
			Showing {filteredEndpoints.length} of {endpoints.length} endpoints
		</div>

	{#each filteredEndpoints as endpoint (endpoint.id)}
		{@const Icon = getStatusIcon(endpoint.status)}
		<Card class="hover:bg-muted/50 transition-colors">
			<CardContent class="p-4">
				<div class="flex items-start gap-4">
					<div class="p-2 border rounded-lg {getStatusBgColor(endpoint.status)}">
						<Icon class="h-5 w-5 {getStatusColor(endpoint.status)}" />
					</div>

						<div class="flex-1 min-w-0">
							<div class="flex items-start justify-between gap-4 mb-3">
								<div class="flex-1">
									<div class="flex items-center gap-2 mb-1">
										<span class="font-semibold text-lg">{endpoint.name}</span>
										<Badge variant="outline" class="text-xs uppercase">{endpoint.type}</Badge>
										<Badge variant={endpoint.status === 'up' ? 'default' : endpoint.status === 'down' ? 'destructive' : 'secondary'} class="capitalize">
											{endpoint.status}
										</Badge>
									</div>
									<div class="flex items-center gap-1 text-sm text-muted-foreground mb-2">
										<Globe class="h-3 w-3" />
										<span class="font-mono">{endpoint.url}</span>
									</div>
								</div>

								<Button
									variant="ghost"
									size="sm"
									onclick={() => toggleExpanded(endpoint.id)}
								>
									{#if expandedEndpoint === endpoint.id}
										<ChevronUp class="h-4 w-4" />
									{:else}
										<ChevronDown class="h-4 w-4" />
									{/if}
								</Button>
							</div>

							<div class="grid grid-cols-4 gap-4 mb-3">
								<div class="flex items-center gap-2">
									<TrendingUp class="h-4 w-4 text-muted-foreground" />
									<div>
										<div class="text-xs text-muted-foreground">Uptime</div>
										<div class="font-semibold" class:text-green-500={endpoint.uptime >= 99.9} class:text-yellow-500={endpoint.uptime >= 95 && endpoint.uptime < 99.9} class:text-red-500={endpoint.uptime < 95}>
											{endpoint.uptime.toFixed(2)}%
										</div>
									</div>
								</div>

								<div class="flex items-center gap-2">
									<Clock class="h-4 w-4 text-muted-foreground" />
									<div>
										<div class="text-xs text-muted-foreground">Response Time</div>
										<div class="font-semibold" class:text-green-500={endpoint.responseTime < 200} class:text-yellow-500={endpoint.responseTime >= 200 && endpoint.responseTime < 1000} class:text-red-500={endpoint.responseTime >= 1000}>
											{formatResponseTime(endpoint.responseTime)}
										</div>
									</div>
								</div>

								<div class="flex items-center gap-2">
									<Activity class="h-4 w-4 text-muted-foreground" />
									<div>
										<div class="text-xs text-muted-foreground">Check Interval</div>
										<div class="font-semibold">{endpoint.checkInterval}s</div>
									</div>
								</div>

								<div class="flex items-center gap-2">
									<AlertTriangle class="h-4 w-4 text-muted-foreground" />
									<div>
										<div class="text-xs text-muted-foreground">Incidents (30d)</div>
										<div class="font-semibold" class:text-green-500={endpoint.incidents.length === 0} class:text-red-500={endpoint.incidents.length > 0}>
											{endpoint.incidents.length}
										</div>
									</div>
								</div>
							</div>

							<div class="text-xs text-muted-foreground">
								Last checked: {formatTimestamp(endpoint.lastChecked)} ({getRelativeTime(endpoint.lastChecked)})
							</div>

							{#if expandedEndpoint === endpoint.id && endpoint.incidents.length > 0}
								<div class="mt-4 pt-4 border-t">
									<div class="text-sm font-semibold mb-3">Recent Incidents</div>
									<div class="space-y-2">
										{#each endpoint.incidents as incident}
											<div class="p-3 bg-muted rounded-lg">
												<div class="flex items-start justify-between mb-1">
													<div class="flex items-center gap-2">
														<AlertTriangle class="h-4 w-4 text-red-500" />
														<span class="font-semibold text-sm">{formatTimestamp(incident.timestamp)}</span>
													</div>
													<Badge variant={incident.duration === null ? 'destructive' : 'outline'}>
														{formatDuration(incident.duration)}
													</Badge>
												</div>
												<p class="text-sm text-muted-foreground ml-6">{incident.reason}</p>
											</div>
										{/each}
									</div>
								</div>
							{/if}

							{#if expandedEndpoint === endpoint.id && endpoint.incidents.length === 0}
								<div class="mt-4 pt-4 border-t">
									<div class="text-sm text-center text-muted-foreground py-4">
										No incidents in the last 30 days
									</div>
								</div>
							{/if}
						</div>
					</div>
				</CardContent>
			</Card>
		{/each}

		{#if filteredEndpoints.length === 0}
			<Card>
				<CardContent class="py-12">
					<div class="text-center text-muted-foreground">
						No endpoints found matching your criteria
					</div>
				</CardContent>
			</Card>
		{/if}
	</div>
</div>
