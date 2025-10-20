<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import { Search, RefreshCw, CheckCircle, XCircle, AlertTriangle, Activity, Cpu, HardDrive, TrendingUp, Clock, GitBranch, ChevronDown, ChevronUp } from 'lucide-svelte';

	let searchQuery = $state('');
	let statusFilter = $state('all');
	let typeFilter = $state('all');
	let isRefreshing = $state(false);
	let expandedService = $state<string | null>(null);

	const statusFilters = [
		{ value: 'all', label: 'All Status' },
		{ value: 'healthy', label: 'Healthy' },
		{ value: 'degraded', label: 'Degraded' },
		{ value: 'down', label: 'Down' }
	];

	const typeFilters = [
		{ value: 'all', label: 'All Types' },
		{ value: 'core', label: 'Core Services' },
		{ value: 'infrastructure', label: 'Infrastructure' },
		{ value: 'worker', label: 'Worker Services' }
	];

	const services = [
		{
			id: 'svc_001',
			name: 'API Gateway',
			type: 'core',
			status: 'healthy',
			uptime: 99.98,
			uptimeSeconds: 2592000,
			version: 'v2.4.1',
			instances: 3,
			activeInstances: 3,
			cpu: 45,
			memory: 68,
			requestsPerMin: 15420,
			avgResponseTime: 89,
			lastRestart: '2024-09-10T10:00:00Z',
			restartCount: 2,
			dependencies: ['PostgreSQL', 'Redis', 'Auth Service'],
			healthChecks: [
				{ name: 'HTTP Health', status: 'passing', lastCheck: '2024-10-10T14:35:00Z' },
				{ name: 'Database Connection', status: 'passing', lastCheck: '2024-10-10T14:35:00Z' },
				{ name: 'Redis Connection', status: 'passing', lastCheck: '2024-10-10T14:35:00Z' }
			]
		},
		{
			id: 'svc_002',
			name: 'Auth Service',
			type: 'core',
			status: 'healthy',
			uptime: 100.0,
			uptimeSeconds: 3456000,
			version: 'v1.8.3',
			instances: 2,
			activeInstances: 2,
			cpu: 28,
			memory: 52,
			requestsPerMin: 3240,
			avgResponseTime: 45,
			lastRestart: '2024-08-15T14:30:00Z',
			restartCount: 1,
			dependencies: ['PostgreSQL', 'Redis'],
			healthChecks: [
				{ name: 'HTTP Health', status: 'passing', lastCheck: '2024-10-10T14:35:05Z' },
				{ name: 'Database Connection', status: 'passing', lastCheck: '2024-10-10T14:35:05Z' }
			]
		},
		{
			id: 'svc_003',
			name: 'PostgreSQL',
			type: 'infrastructure',
			status: 'healthy',
			uptime: 99.95,
			uptimeSeconds: 4320000,
			version: 'v16.2',
			instances: 3,
			activeInstances: 3,
			cpu: 65,
			memory: 82,
			requestsPerMin: 42000,
			avgResponseTime: 8,
			lastRestart: '2024-07-20T08:00:00Z',
			restartCount: 3,
			dependencies: [],
			healthChecks: [
				{ name: 'Primary Health', status: 'passing', lastCheck: '2024-10-10T14:35:10Z' },
				{ name: 'Replication Status', status: 'passing', lastCheck: '2024-10-10T14:35:10Z' },
				{ name: 'Connection Pool', status: 'passing', lastCheck: '2024-10-10T14:35:10Z' }
			]
		},
		{
			id: 'svc_004',
			name: 'Redis Cache',
			type: 'infrastructure',
			status: 'healthy',
			uptime: 100.0,
			uptimeSeconds: 5184000,
			version: 'v7.2.4',
			instances: 2,
			activeInstances: 2,
			cpu: 18,
			memory: 94,
			requestsPerMin: 89000,
			avgResponseTime: 2,
			lastRestart: '2024-06-01T12:00:00Z',
			restartCount: 0,
			dependencies: [],
			healthChecks: [
				{ name: 'Master Health', status: 'passing', lastCheck: '2024-10-10T14:35:15Z' },
				{ name: 'Replica Sync', status: 'passing', lastCheck: '2024-10-10T14:35:15Z' }
			]
		},
		{
			id: 'svc_005',
			name: 'Build Service',
			type: 'worker',
			status: 'degraded',
			uptime: 97.3,
			uptimeSeconds: 1728000,
			version: 'v3.1.0',
			instances: 4,
			activeInstances: 3,
			cpu: 88,
			memory: 76,
			requestsPerMin: 120,
			avgResponseTime: 2340,
			lastRestart: '2024-10-10T12:00:00Z',
			restartCount: 8,
			dependencies: ['Registry', 'Storage Service', 'Queue Service'],
			healthChecks: [
				{ name: 'HTTP Health', status: 'passing', lastCheck: '2024-10-10T14:35:20Z' },
				{ name: 'Build Queue', status: 'warning', lastCheck: '2024-10-10T14:35:20Z' },
				{ name: 'Storage Connection', status: 'passing', lastCheck: '2024-10-10T14:35:20Z' }
			]
		},
		{
			id: 'svc_006',
			name: 'Queue Service',
			type: 'worker',
			status: 'healthy',
			uptime: 99.87,
			uptimeSeconds: 2160000,
			version: 'v2.3.5',
			instances: 3,
			activeInstances: 3,
			cpu: 42,
			memory: 58,
			requestsPerMin: 4560,
			avgResponseTime: 125,
			lastRestart: '2024-09-15T16:00:00Z',
			restartCount: 4,
			dependencies: ['Redis Cache', 'PostgreSQL'],
			healthChecks: [
				{ name: 'HTTP Health', status: 'passing', lastCheck: '2024-10-10T14:35:25Z' },
				{ name: 'Queue Processing', status: 'passing', lastCheck: '2024-10-10T14:35:25Z' },
				{ name: 'Redis Connection', status: 'passing', lastCheck: '2024-10-10T14:35:25Z' }
			]
		},
		{
			id: 'svc_007',
			name: 'Registry',
			type: 'infrastructure',
			status: 'healthy',
			uptime: 99.92,
			uptimeSeconds: 3024000,
			version: 'v2.8.2',
			instances: 2,
			activeInstances: 2,
			cpu: 38,
			memory: 48,
			requestsPerMin: 2340,
			avgResponseTime: 156,
			lastRestart: '2024-08-20T10:00:00Z',
			restartCount: 2,
			dependencies: ['Storage Service'],
			healthChecks: [
				{ name: 'HTTP Health', status: 'passing', lastCheck: '2024-10-10T14:35:30Z' },
				{ name: 'Storage Connection', status: 'passing', lastCheck: '2024-10-10T14:35:30Z' }
			]
		},
		{
			id: 'svc_008',
			name: 'Scheduler',
			type: 'worker',
			status: 'healthy',
			uptime: 99.96,
			uptimeSeconds: 2592000,
			version: 'v1.5.2',
			instances: 2,
			activeInstances: 2,
			cpu: 12,
			memory: 32,
			requestsPerMin: 180,
			avgResponseTime: 45,
			lastRestart: '2024-09-10T08:00:00Z',
			restartCount: 1,
			dependencies: ['PostgreSQL', 'Queue Service'],
			healthChecks: [
				{ name: 'HTTP Health', status: 'passing', lastCheck: '2024-10-10T14:35:35Z' },
				{ name: 'Cron Jobs', status: 'passing', lastCheck: '2024-10-10T14:35:35Z' }
			]
		},
		{
			id: 'svc_009',
			name: 'Webhook Service',
			type: 'worker',
			status: 'healthy',
			uptime: 99.82,
			uptimeSeconds: 1987200,
			version: 'v2.1.4',
			instances: 3,
			activeInstances: 3,
			cpu: 35,
			memory: 45,
			requestsPerMin: 890,
			avgResponseTime: 234,
			lastRestart: '2024-09-25T12:00:00Z',
			restartCount: 5,
			dependencies: ['Queue Service', 'API Gateway'],
			healthChecks: [
				{ name: 'HTTP Health', status: 'passing', lastCheck: '2024-10-10T14:35:40Z' },
				{ name: 'Webhook Queue', status: 'passing', lastCheck: '2024-10-10T14:35:40Z' }
			]
		},
		{
			id: 'svc_010',
			name: 'Metrics Collector',
			type: 'infrastructure',
			status: 'down',
			uptime: 89.5,
			uptimeSeconds: 0,
			version: 'v1.9.1',
			instances: 2,
			activeInstances: 0,
			cpu: 0,
			memory: 0,
			requestsPerMin: 0,
			avgResponseTime: 0,
			lastRestart: '2024-10-10T10:00:00Z',
			restartCount: 12,
			dependencies: ['PostgreSQL', 'Redis Cache'],
			healthChecks: [
				{ name: 'HTTP Health', status: 'critical', lastCheck: '2024-10-10T14:35:45Z' },
				{ name: 'Database Connection', status: 'critical', lastCheck: '2024-10-10T14:35:45Z' },
				{ name: 'Metrics Export', status: 'critical', lastCheck: '2024-10-10T14:35:45Z' }
			]
		}
	];

	const filteredServices = $derived(
		services.filter((service) => {
			const matchesSearch =
				service.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				service.version.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesStatus = statusFilter === 'all' || service.status === statusFilter;
			const matchesType = typeFilter === 'all' || service.type === typeFilter;
			return matchesSearch && matchesStatus && matchesType;
		})
	);

	const stats = $derived({
		total: services.length,
		healthy: services.filter((s) => s.status === 'healthy').length,
		degraded: services.filter((s) => s.status === 'degraded').length,
		down: services.filter((s) => s.status === 'down').length,
		avgUptime: (services.reduce((sum, s) => sum + s.uptime, 0) / services.length).toFixed(2),
		totalInstances: services.reduce((sum, s) => sum + s.instances, 0),
		activeInstances: services.reduce((sum, s) => sum + s.activeInstances, 0)
	});

	function getStatusColor(status: string) {
		switch (status) {
			case 'healthy':
				return 'text-green-500';
			case 'degraded':
				return 'text-yellow-500';
			case 'down':
				return 'text-red-500';
			default:
				return 'text-gray-500';
		}
	}

	function getStatusBgColor(status: string) {
		switch (status) {
			case 'healthy':
				return 'bg-green-500/10 border-green-500/20';
			case 'degraded':
				return 'bg-yellow-500/10 border-yellow-500/20';
			case 'down':
				return 'bg-red-500/10 border-red-500/20';
			default:
				return 'bg-gray-500/10 border-gray-500/20';
		}
	}

	function getStatusIcon(status: string) {
		switch (status) {
			case 'healthy':
				return CheckCircle;
			case 'degraded':
				return AlertTriangle;
			case 'down':
				return XCircle;
			default:
				return Activity;
		}
	}

	function getHealthCheckIcon(status: string) {
		switch (status) {
			case 'passing':
				return CheckCircle;
			case 'warning':
				return AlertTriangle;
			case 'critical':
				return XCircle;
			default:
				return Activity;
		}
	}

	function getHealthCheckColor(status: string) {
		switch (status) {
			case 'passing':
				return 'text-green-500';
			case 'warning':
				return 'text-yellow-500';
			case 'critical':
				return 'text-red-500';
			default:
				return 'text-gray-500';
		}
	}

	function formatNumber(num: number) {
		if (num >= 1000000) return `${(num / 1000000).toFixed(2)}M`;
		if (num >= 1000) return `${(num / 1000).toFixed(1)}K`;
		return num.toString();
	}

	function formatUptime(seconds: number) {
		if (seconds === 0) return 'Down';
		const days = Math.floor(seconds / 86400);
		const hours = Math.floor((seconds % 86400) / 3600);
		if (days > 0) return `${days}d ${hours}h`;
		return `${hours}h`;
	}

	function formatResponseTime(ms: number) {
		if (ms === 0) return 'N/A';
		if (ms < 1000) return `${ms}ms`;
		return `${(ms / 1000).toFixed(2)}s`;
	}

	function getResourceColor(usage: number) {
		if (usage >= 90) return 'text-red-500';
		if (usage >= 75) return 'text-yellow-500';
		return 'text-green-500';
	}

	function getResourceBgColor(usage: number) {
		if (usage >= 90) return 'bg-red-500';
		if (usage >= 75) return 'bg-yellow-500';
		return 'bg-green-500';
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
		expandedService = expandedService === id ? null : id;
	}

	async function handleRefresh() {
		isRefreshing = true;
		await new Promise((resolve) => setTimeout(resolve, 1000));
		isRefreshing = false;
	}

	function exportServices() {
		const dataStr = JSON.stringify(filteredServices, null, 2);
		const blob = new Blob([dataStr], { type: 'application/json' });
		const url = URL.createObjectURL(blob);
		const link = document.createElement('a');
		link.href = url;
		link.download = `services-${new Date().toISOString()}.json`;
		link.click();
		URL.revokeObjectURL(url);
	}
</script>

<div class="flex-1 overflow-y-auto">
	<div class="border-b border-border bg-background sticky top-0 z-10">
		<div class="px-8 py-6">
			<div class="flex items-center justify-between mb-4">
				<div>
					<h1 class="text-3xl font-bold">Services</h1>
					<p class="text-muted-foreground mt-1">Monitor internal PaaS services and infrastructure</p>
				</div>
				<div class="flex items-center gap-2">
					<Button variant="outline" size="sm" onclick={exportServices}>
						Export
					</Button>
				<Button variant="outline" size="icon" onclick={handleRefresh} disabled={isRefreshing}>
					<RefreshCw class={`h-4 w-4 ${isRefreshing ? 'animate-spin' : ''}`} />
				</Button>
				</div>
			</div>

			<div class="grid grid-cols-7 gap-4 mb-4">
				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Total Services</CardDescription>
						<CardTitle class="text-2xl">{stats.total}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Healthy</CardDescription>
						<CardTitle class="text-2xl text-green-500">{stats.healthy}</CardTitle>
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
						<CardDescription class="text-xs">Down</CardDescription>
						<CardTitle class="text-2xl text-red-500">{stats.down}</CardTitle>
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
						<CardDescription class="text-xs">Instances</CardDescription>
						<CardTitle class="text-2xl text-purple-500">{stats.activeInstances}/{stats.totalInstances}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Health Score</CardDescription>
						<CardTitle class="text-2xl text-green-500">{Math.round((stats.healthy / stats.total) * 100)}%</CardTitle>
					</CardHeader>
				</Card>
			</div>

			<div class="grid grid-cols-3 gap-3">
				<div class="relative">
					<Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
					<Input
						type="text"
						placeholder="Search services..."
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
			Showing {filteredServices.length} of {services.length} services
		</div>

	{#each filteredServices as service (service.id)}
		{@const Icon = getStatusIcon(service.status)}
		<Card class="hover:bg-muted/50 transition-colors">
			<CardContent class="p-4">
				<div class="flex items-start gap-4">
					<div class="p-2 border rounded-lg {getStatusBgColor(service.status)}">
						<Icon class="h-5 w-5 {getStatusColor(service.status)}" />
					</div>

						<div class="flex-1 min-w-0">
							<div class="flex items-start justify-between gap-4 mb-3">
								<div class="flex-1">
									<div class="flex items-center gap-2 mb-1">
										<span class="font-semibold text-lg">{service.name}</span>
										<Badge variant="outline" class="text-xs font-mono">{service.version}</Badge>
										<Badge variant={service.status === 'healthy' ? 'default' : service.status === 'down' ? 'destructive' : 'secondary'} class="capitalize">
											{service.status}
										</Badge>
										<Badge variant="outline" class="text-xs capitalize">{service.type}</Badge>
									</div>
									<div class="flex items-center gap-3 text-sm text-muted-foreground mb-2">
										<span>Instances: {service.activeInstances}/{service.instances}</span>
										<span>•</span>
										<span>Uptime: {service.uptime.toFixed(2)}%</span>
										<span>•</span>
										<span>Restarts: {service.restartCount}</span>
									</div>
								</div>

								<Button
									variant="ghost"
									size="sm"
									onclick={() => toggleExpanded(service.id)}
								>
									{#if expandedService === service.id}
										<ChevronUp class="h-4 w-4" />
									{:else}
										<ChevronDown class="h-4 w-4" />
									{/if}
								</Button>
							</div>

							<div class="grid grid-cols-4 gap-4 mb-3">
								<div>
									<div class="flex items-center gap-2 mb-2">
										<Cpu class="h-4 w-4 text-muted-foreground" />
										<span class="text-xs text-muted-foreground">CPU Usage</span>
									</div>
									<div class="space-y-1">
										<div class="flex items-center justify-between text-sm">
											<span class="font-semibold {getResourceColor(service.cpu)}">{service.cpu}%</span>
										</div>
										<div class="h-2 bg-muted rounded-full overflow-hidden">
											<div class="{getResourceBgColor(service.cpu)} h-full transition-all" style="width: {service.cpu}%"></div>
										</div>
									</div>
								</div>

								<div>
									<div class="flex items-center gap-2 mb-2">
										<HardDrive class="h-4 w-4 text-muted-foreground" />
										<span class="text-xs text-muted-foreground">Memory Usage</span>
									</div>
									<div class="space-y-1">
										<div class="flex items-center justify-between text-sm">
											<span class="font-semibold {getResourceColor(service.memory)}">{service.memory}%</span>
										</div>
										<div class="h-2 bg-muted rounded-full overflow-hidden">
											<div class="{getResourceBgColor(service.memory)} h-full transition-all" style="width: {service.memory}%"></div>
										</div>
									</div>
								</div>

								<div class="flex items-center gap-2">
									<TrendingUp class="h-4 w-4 text-muted-foreground" />
									<div>
										<div class="text-xs text-muted-foreground">Requests/min</div>
										<div class="font-semibold">{formatNumber(service.requestsPerMin)}</div>
									</div>
								</div>

								<div class="flex items-center gap-2">
									<Clock class="h-4 w-4 text-muted-foreground" />
									<div>
										<div class="text-xs text-muted-foreground">Avg Response</div>
										<div class="font-semibold" class:text-green-500={service.avgResponseTime < 200} class:text-yellow-500={service.avgResponseTime >= 200 && service.avgResponseTime < 1000} class:text-red-500={service.avgResponseTime >= 1000}>
											{formatResponseTime(service.avgResponseTime)}
										</div>
									</div>
								</div>
							</div>

							<div class="text-xs text-muted-foreground">
								Last restart: {formatTimestamp(service.lastRestart)} ({getRelativeTime(service.lastRestart)})
							</div>

							{#if expandedService === service.id}
								<div class="mt-4 pt-4 border-t space-y-4">
									<div>
										<div class="text-sm font-semibold mb-2">Health Checks</div>
										<div class="space-y-2">
									{#each service.healthChecks as check}
										{@const CheckIcon = getHealthCheckIcon(check.status)}
										<div class="flex items-center justify-between p-2 bg-muted rounded-lg">
											<div class="flex items-center gap-2">
												<CheckIcon class="h-4 w-4 {getHealthCheckColor(check.status)}" />
												<span class="text-sm font-medium">{check.name}</span>
											</div>
													<div class="flex items-center gap-2">
														<Badge variant={check.status === 'passing' ? 'default' : check.status === 'critical' ? 'destructive' : 'secondary'} class="capitalize text-xs">
															{check.status}
														</Badge>
														<span class="text-xs text-muted-foreground">{getRelativeTime(check.lastCheck)}</span>
													</div>
												</div>
											{/each}
										</div>
									</div>

									{#if service.dependencies.length > 0}
										<div>
											<div class="text-sm font-semibold mb-2 flex items-center gap-2">
												<GitBranch class="h-4 w-4" />
												<span>Dependencies</span>
											</div>
											<div class="flex flex-wrap gap-2">
												{#each service.dependencies as dep}
													<Badge variant="outline" class="text-xs">{dep}</Badge>
												{/each}
											</div>
										</div>
									{/if}

									<div class="grid grid-cols-2 gap-4">
										<div class="space-y-2">
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">Service Name:</span>
												<span class="font-medium">{service.name}</span>
											</div>
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">Version:</span>
												<span class="font-mono">{service.version}</span>
											</div>
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">Type:</span>
												<Badge variant="outline" class="text-xs capitalize">{service.type}</Badge>
											</div>
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">Uptime:</span>
												<span>{formatUptime(service.uptimeSeconds)}</span>
											</div>
										</div>
										<div class="space-y-2">
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">Total Instances:</span>
												<span>{service.instances}</span>
											</div>
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">Active Instances:</span>
												<span class={service.activeInstances === service.instances ? 'text-green-500' : 'text-yellow-500'}>{service.activeInstances}</span>
											</div>
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">Restart Count:</span>
												<span class={service.restartCount > 5 ? 'text-red-500' : service.restartCount > 2 ? 'text-yellow-500' : ''}>{service.restartCount}</span>
											</div>
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">Availability:</span>
												<span class={service.uptime >= 99.9 ? 'text-green-500' : service.uptime >= 99 ? 'text-yellow-500' : 'text-red-500'}>{service.uptime.toFixed(2)}%</span>
											</div>
										</div>
									</div>
								</div>
							{/if}
						</div>
					</div>
				</CardContent>
			</Card>
		{/each}

		{#if filteredServices.length === 0}
			<Card>
				<CardContent class="py-12">
					<div class="text-center text-muted-foreground">
						No services found matching your criteria
					</div>
				</CardContent>
			</Card>
		{/if}
	</div>
</div>
