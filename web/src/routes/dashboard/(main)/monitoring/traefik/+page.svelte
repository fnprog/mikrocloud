<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import { Search, RefreshCw, Activity, TrendingUp, Globe, Zap, Shield, AlertCircle } from 'lucide-svelte';

	let searchQuery = $state('');
	let statusFilter = $state('all');
	let typeFilter = $state('all');
	let isRefreshing = $state(false);

	const statusFilters = [
		{ value: 'all', label: 'All Status' },
		{ value: 'healthy', label: 'Healthy' },
		{ value: 'warning', label: 'Warning' },
		{ value: 'error', label: 'Error' }
	];

	const typeFilters = [
		{ value: 'all', label: 'All Types' },
		{ value: 'http', label: 'HTTP' },
		{ value: 'tcp', label: 'TCP' },
		{ value: 'udp', label: 'UDP' }
	];

	const routes = [
		{
			id: 'route_001',
			name: 'mikrocloud-app',
			rule: 'Host(`app.mikrocloud.io`)',
			entrypoint: 'websecure',
			service: 'mikrocloud-app-service',
			type: 'http',
			status: 'healthy',
			requests: 125430,
			requestRate: 45.2,
			avgResponseTime: 142,
			statusCodes: {
				'2xx': 122350,
				'3xx': 2180,
				'4xx': 850,
				'5xx': 50
			},
			tls: true,
			middleware: ['compress', 'rate-limit', 'headers']
		},
		{
			id: 'route_002',
			name: 'api-gateway',
			rule: 'Host(`api.mikrocloud.io`) && PathPrefix(`/v1`)',
			entrypoint: 'websecure',
			service: 'api-gateway-service',
			type: 'http',
			status: 'healthy',
			requests: 458920,
			requestRate: 128.7,
			avgResponseTime: 89,
			statusCodes: {
				'2xx': 445200,
				'3xx': 5120,
				'4xx': 7800,
				'5xx': 800
			},
			tls: true,
			middleware: ['auth', 'cors', 'rate-limit', 'compress']
		},
		{
			id: 'route_003',
			name: 'staging-env',
			rule: 'Host(`staging.mikrocloud.io`)',
			entrypoint: 'websecure',
			service: 'staging-service',
			type: 'http',
			status: 'warning',
			requests: 8450,
			requestRate: 2.8,
			avgResponseTime: 2340,
			statusCodes: {
				'2xx': 6200,
				'3xx': 180,
				'4xx': 1570,
				'5xx': 500
			},
			tls: true,
			middleware: ['compress', 'headers']
		},
		{
			id: 'route_004',
			name: 'cdn-edge',
			rule: 'Host(`cdn.mikrocloud.io`)',
			entrypoint: 'websecure',
			service: 'cdn-service',
			type: 'http',
			status: 'healthy',
			requests: 1245800,
			requestRate: 342.5,
			avgResponseTime: 28,
			statusCodes: {
				'2xx': 1238900,
				'3xx': 4200,
				'4xx': 2500,
				'5xx': 200
			},
			tls: true,
			middleware: ['cache', 'compress']
		},
		{
			id: 'route_005',
			name: 'webhooks',
			rule: 'Host(`webhooks.mikrocloud.io`)',
			entrypoint: 'websecure',
			service: 'webhook-service',
			type: 'http',
			status: 'healthy',
			requests: 24580,
			requestRate: 8.4,
			avgResponseTime: 156,
			statusCodes: {
				'2xx': 23200,
				'3xx': 420,
				'4xx': 860,
				'5xx': 100
			},
			tls: true,
			middleware: ['auth', 'rate-limit']
		},
		{
			id: 'route_006',
			name: 'legacy-api',
			rule: 'Host(`legacy.internal`) && PathPrefix(`/api`)',
			entrypoint: 'web',
			service: 'legacy-api-service',
			type: 'http',
			status: 'error',
			requests: 1240,
			requestRate: 0.2,
			avgResponseTime: 0,
			statusCodes: {
				'2xx': 0,
				'3xx': 0,
				'4xx': 0,
				'5xx': 1240
			},
			tls: false,
			middleware: []
		},
		{
			id: 'route_007',
			name: 'postgres-lb',
			rule: 'HostSNI(`*`)',
			entrypoint: 'postgres-tcp',
			service: 'postgres-cluster',
			type: 'tcp',
			status: 'healthy',
			requests: 58920,
			requestRate: 18.5,
			avgResponseTime: 12,
			statusCodes: {},
			tls: false,
			middleware: []
		},
		{
			id: 'route_008',
			name: 'redis-proxy',
			rule: 'HostSNI(`*`)',
			entrypoint: 'redis-tcp',
			service: 'redis-cluster',
			type: 'tcp',
			status: 'healthy',
			requests: 245800,
			requestRate: 82.4,
			avgResponseTime: 5,
			statusCodes: {},
			tls: false,
			middleware: []
		},
		{
			id: 'route_009',
			name: 'metrics-collector',
			rule: 'Host(`metrics.internal`)',
			entrypoint: 'web',
			service: 'metrics-service',
			type: 'http',
			status: 'healthy',
			requests: 12450,
			requestRate: 4.2,
			avgResponseTime: 67,
			statusCodes: {
				'2xx': 12380,
				'3xx': 20,
				'4xx': 40,
				'5xx': 10
			},
			tls: false,
			middleware: ['compress']
		},
		{
			id: 'route_010',
			name: 'build-server',
			rule: 'Host(`builds.mikrocloud.io`)',
			entrypoint: 'websecure',
			service: 'build-service',
			type: 'http',
			status: 'warning',
			requests: 3580,
			requestRate: 1.2,
			avgResponseTime: 1850,
			statusCodes: {
				'2xx': 2900,
				'3xx': 120,
				'4xx': 480,
				'5xx': 80
			},
			tls: true,
			middleware: ['auth', 'compress', 'timeout']
		}
	];

	const filteredRoutes = $derived(
		routes.filter((route) => {
			const matchesSearch =
				route.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				route.rule.toLowerCase().includes(searchQuery.toLowerCase()) ||
				route.service.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesStatus = statusFilter === 'all' || route.status === statusFilter;
			const matchesType = typeFilter === 'all' || route.type === typeFilter;
			return matchesSearch && matchesStatus && matchesType;
		})
	);

	const stats = $derived({
		totalRoutes: routes.length,
		totalRequests: routes.reduce((sum, r) => sum + r.requests, 0),
		totalRequestRate: routes.reduce((sum, r) => sum + r.requestRate, 0).toFixed(1),
		avgResponseTime: Math.round(
			routes.filter((r) => r.status !== 'error').reduce((sum, r) => sum + r.avgResponseTime, 0) /
				routes.filter((r) => r.status !== 'error').length
		),
		healthy: routes.filter((r) => r.status === 'healthy').length,
		warning: routes.filter((r) => r.status === 'warning').length,
		error: routes.filter((r) => r.status === 'error').length
	});

	function getStatusColor(status: string) {
		switch (status) {
			case 'healthy':
				return 'text-green-500';
			case 'warning':
				return 'text-yellow-500';
			case 'error':
				return 'text-red-500';
			default:
				return 'text-gray-500';
		}
	}

	function getStatusBgColor(status: string) {
		switch (status) {
			case 'healthy':
				return 'bg-green-500/10 border-green-500/20';
			case 'warning':
				return 'bg-yellow-500/10 border-yellow-500/20';
			case 'error':
				return 'bg-red-500/10 border-red-500/20';
			default:
				return 'bg-gray-500/10 border-gray-500/20';
		}
	}

	function formatNumber(num: number) {
		if (num >= 1000000) return `${(num / 1000000).toFixed(2)}M`;
		if (num >= 1000) return `${(num / 1000).toFixed(1)}K`;
		return num.toString();
	}

	function formatResponseTime(ms: number) {
		if (ms === 0) return 'N/A';
		if (ms < 1000) return `${ms}ms`;
		return `${(ms / 1000).toFixed(2)}s`;
	}

	function calculateSuccessRate(statusCodes: Record<string, number>) {
		const total = Object.values(statusCodes).reduce((sum, count) => sum + count, 0);
		if (total === 0) return 0;
		const success = (statusCodes['2xx'] || 0) + (statusCodes['3xx'] || 0);
		return ((success / total) * 100).toFixed(2);
	}

	async function handleRefresh() {
		isRefreshing = true;
		await new Promise((resolve) => setTimeout(resolve, 1000));
		isRefreshing = false;
	}

	function exportRoutes() {
		const dataStr = JSON.stringify(filteredRoutes, null, 2);
		const blob = new Blob([dataStr], { type: 'application/json' });
		const url = URL.createObjectURL(blob);
		const link = document.createElement('a');
		link.href = url;
		link.download = `traefik-routes-${new Date().toISOString()}.json`;
		link.click();
		URL.revokeObjectURL(url);
	}
</script>

<div class="flex-1 overflow-y-auto">
	<div class="border-b border-border bg-background sticky top-0 z-10">
		<div class="px-8 py-6">
			<div class="flex items-center justify-between mb-4">
				<div>
					<h1 class="text-3xl font-bold">Traefik Proxy</h1>
					<p class="text-muted-foreground mt-1">Monitor routing metrics and traffic</p>
				</div>
				<div class="flex items-center gap-2">
					<Button variant="outline" size="sm" onclick={exportRoutes}>
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
						<CardDescription class="text-xs">Total Routes</CardDescription>
						<CardTitle class="text-2xl">{stats.totalRoutes}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Total Requests</CardDescription>
						<CardTitle class="text-2xl text-blue-500">{formatNumber(stats.totalRequests)}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Request Rate</CardDescription>
						<CardTitle class="text-2xl text-purple-500">{stats.totalRequestRate}/s</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Avg Response</CardDescription>
						<CardTitle class="text-2xl text-cyan-500">{stats.avgResponseTime}ms</CardTitle>
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
						<CardDescription class="text-xs">Warning</CardDescription>
						<CardTitle class="text-2xl text-yellow-500">{stats.warning}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Error</CardDescription>
						<CardTitle class="text-2xl text-red-500">{stats.error}</CardTitle>
					</CardHeader>
				</Card>
			</div>

			<div class="grid grid-cols-3 gap-3">
				<div class="relative">
					<Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
					<Input
						type="text"
						placeholder="Search routes..."
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
			Showing {filteredRoutes.length} of {routes.length} routes
		</div>

		{#each filteredRoutes as route (route.id)}
			<Card class="hover:bg-muted/50 transition-colors">
				<CardContent class="p-4">
					<div class="flex items-start gap-4">
						<div class="p-2 border rounded-lg {getStatusBgColor(route.status)}">
							<Activity class="h-5 w-5 {getStatusColor(route.status)}" />
						</div>

						<div class="flex-1 min-w-0">
							<div class="flex items-start justify-between gap-4 mb-3">
								<div class="flex-1">
									<div class="flex items-center gap-2 mb-1">
										<span class="font-semibold text-lg">{route.name}</span>
										<Badge variant="outline" class="text-xs uppercase">{route.type}</Badge>
										<Badge variant={route.status === 'healthy' ? 'default' : route.status === 'error' ? 'destructive' : 'secondary'} class="capitalize">
											{route.status}
										</Badge>
										{#if route.tls}
											<Badge variant="outline" class="text-xs">
												<Shield class="h-3 w-3 mr-1" />
												TLS
											</Badge>
										{/if}
									</div>
									<div class="flex items-center gap-1 text-sm text-muted-foreground mb-1">
										<Globe class="h-3 w-3" />
										<span class="font-mono text-xs">{route.rule}</span>
									</div>
									<div class="flex items-center gap-2 text-xs text-muted-foreground">
										<span>Entrypoint: <span class="font-semibold">{route.entrypoint}</span></span>
										<span>·</span>
										<span>Service: <span class="font-semibold">{route.service}</span></span>
									</div>
								</div>
							</div>

							<div class="grid grid-cols-4 gap-4 mb-3">
								<div class="flex items-center gap-2">
									<Activity class="h-4 w-4 text-muted-foreground" />
									<div>
										<div class="text-xs text-muted-foreground">Requests</div>
										<div class="font-semibold">{formatNumber(route.requests)}</div>
									</div>
								</div>

								<div class="flex items-center gap-2">
									<Zap class="h-4 w-4 text-muted-foreground" />
									<div>
										<div class="text-xs text-muted-foreground">Request Rate</div>
										<div class="font-semibold">{route.requestRate.toFixed(1)}/s</div>
									</div>
								</div>

								<div class="flex items-center gap-2">
									<TrendingUp class="h-4 w-4 text-muted-foreground" />
									<div>
										<div class="text-xs text-muted-foreground">Response Time</div>
										<div class="font-semibold" class:text-green-500={route.avgResponseTime < 200} class:text-yellow-500={route.avgResponseTime >= 200 && route.avgResponseTime < 1000} class:text-red-500={route.avgResponseTime >= 1000}>
											{formatResponseTime(route.avgResponseTime)}
										</div>
									</div>
								</div>

								{#if route.type === 'http'}
									<div class="flex items-center gap-2">
										<Activity class="h-4 w-4 text-muted-foreground" />
										<div>
											<div class="text-xs text-muted-foreground">Success Rate</div>
											<div class="font-semibold" class:text-green-500={parseFloat(calculateSuccessRate(route.statusCodes)) >= 99} class:text-yellow-500={parseFloat(calculateSuccessRate(route.statusCodes)) >= 95 && parseFloat(calculateSuccessRate(route.statusCodes)) < 99} class:text-red-500={parseFloat(calculateSuccessRate(route.statusCodes)) < 95}>
												{calculateSuccessRate(route.statusCodes)}%
											</div>
										</div>
									</div>
								{/if}
							</div>

							{#if route.type === 'http' && Object.keys(route.statusCodes).length > 0}
								<div class="mb-3 p-3 bg-muted/50 rounded-lg">
									<div class="text-xs font-semibold mb-2">HTTP Status Codes</div>
									<div class="grid grid-cols-4 gap-3">
										<div>
											<div class="text-xs text-muted-foreground">2xx (Success)</div>
											<div class="font-semibold text-green-500">{formatNumber(route.statusCodes['2xx'] || 0)}</div>
										</div>
										<div>
											<div class="text-xs text-muted-foreground">3xx (Redirect)</div>
											<div class="font-semibold text-blue-500">{formatNumber(route.statusCodes['3xx'] || 0)}</div>
										</div>
										<div>
											<div class="text-xs text-muted-foreground">4xx (Client Error)</div>
											<div class="font-semibold text-yellow-500">{formatNumber(route.statusCodes['4xx'] || 0)}</div>
										</div>
										<div>
											<div class="text-xs text-muted-foreground">5xx (Server Error)</div>
											<div class="font-semibold text-red-500">{formatNumber(route.statusCodes['5xx'] || 0)}</div>
										</div>
									</div>
								</div>
							{/if}

							{#if route.middleware.length > 0}
								<div class="flex items-center gap-2 text-xs">
									<span class="text-muted-foreground">Middleware:</span>
									<div class="flex gap-1 flex-wrap">
										{#each route.middleware as mw}
											<Badge variant="outline" class="text-xs">{mw}</Badge>
										{/each}
									</div>
								</div>
							{/if}

							{#if route.status === 'error'}
								<div class="mt-3 p-2 bg-red-500/10 border border-red-500/20 rounded flex items-start gap-2">
									<AlertCircle class="h-4 w-4 text-red-500 mt-0.5" />
									<div class="text-sm text-red-500">
										<span class="font-semibold">Service Unavailable:</span> All requests are failing (100% 5xx responses)
									</div>
								</div>
							{/if}

							{#if route.status === 'warning'}
								<div class="mt-3 p-2 bg-yellow-500/10 border border-yellow-500/20 rounded flex items-start gap-2">
									<AlertCircle class="h-4 w-4 text-yellow-500 mt-0.5" />
									<div class="text-sm text-yellow-500">
										<span class="font-semibold">Performance Degraded:</span> High response times or elevated error rate detected
									</div>
								</div>
							{/if}
						</div>
					</div>
				</CardContent>
			</Card>
		{/each}

		{#if filteredRoutes.length === 0}
			<Card>
				<CardContent class="py-12">
					<div class="text-center text-muted-foreground">
						No routes found matching your criteria
					</div>
				</CardContent>
			</Card>
		{/if}
	</div>
</div>
