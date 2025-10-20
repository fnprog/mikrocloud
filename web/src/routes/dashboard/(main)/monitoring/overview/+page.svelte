<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Select from '$lib/components/ui/select';
	import { Calendar, RefreshCw, MoreVertical } from 'lucide-svelte';

	let selectedEnvironment = $state('production');
	let timeRange = $state('last-12-hours');
	let isRefreshing = $state(false);

	const environments = [
		{ value: 'production', label: 'Production' },
		{ value: 'staging', label: 'Staging' },
		{ value: 'development', label: 'Development' }
	];

	const timeRanges = [
		{ value: 'last-hour', label: 'Last hour' },
		{ value: 'last-12-hours', label: 'Last 12 hours' },
		{ value: 'last-24-hours', label: 'Last 24 hours' },
		{ value: 'last-7-days', label: 'Last 7 days' },
		{ value: 'last-30-days', label: 'Last 30 days' }
	];

	const dummyContainerMetrics = [
		{
			id: 'mikrocloud-api',
			name: 'mikrocloud-api',
			cpu: 45.2,
			memory: { used: '512MB', limit: '1GB', percent: 50 },
			network: { in: '1.2 MB/s', out: '800 KB/s' },
			disk: { read: '50 KB/s', write: '25 KB/s' },
			restarts: 0,
			state: 'running',
			uptime: '15d 4h 23m'
		},
		{
			id: 'mikrocloud-scheduler',
			name: 'mikrocloud-scheduler',
			cpu: 12.5,
			memory: { used: '256MB', limit: '512MB', percent: 50 },
			network: { in: '100 KB/s', out: '50 KB/s' },
			disk: { read: '10 KB/s', write: '5 KB/s' },
			restarts: 1,
			state: 'running',
			uptime: '15d 4h 20m'
		},
		{
			id: 'app-web-1',
			name: 'app-web-1',
			cpu: 78.9,
			memory: { used: '1.5GB', limit: '2GB', percent: 75 },
			network: { in: '5.5 MB/s', out: '3.2 MB/s' },
			disk: { read: '200 KB/s', write: '150 KB/s' },
			restarts: 0,
			state: 'running',
			uptime: '7d 12h 45m'
		},
		{
			id: 'postgres-main',
			name: 'postgres-main',
			cpu: 32.1,
			memory: { used: '768MB', limit: '2GB', percent: 38 },
			network: { in: '2.1 MB/s', out: '1.8 MB/s' },
			disk: { read: '500 KB/s', write: '350 KB/s' },
			restarts: 0,
			state: 'running',
			uptime: '30d 8h 15m'
		}
	];

	const serverMetrics = {
		cpuLoadAvg: [1.2, 1.5, 1.8],
		memory: { total: '16GB', used: '8.2GB', free: '7.8GB', percent: 51 },
		disk: { total: '500GB', used: '245GB', free: '255GB', percent: 49 },
		uptime: '45d 12h 30m'
	};

	const traefikMetrics = {
		totalRequests: 1245678,
		requestsPerSec: 234,
		avgLatency: 45,
		p95Latency: 120,
		p99Latency: 250,
		errorRate: 0.5
	};

	const paasServices = [
		{ name: 'API Server', status: 'healthy', requests: 45678, latency: 45, errors: 12 },
		{ name: 'Scheduler', status: 'healthy', requests: 8234, latency: 23, errors: 3 },
		{ name: 'Registry', status: 'healthy', requests: 2145, latency: 67, errors: 0 },
		{ name: 'Build Worker', status: 'degraded', requests: 456, latency: 230, errors: 15 }
	];

	function getStateColor(state: string) {
		switch (state) {
			case 'running':
				return 'bg-green-500';
			case 'exited':
				return 'bg-red-500';
			case 'paused':
				return 'bg-yellow-500';
			default:
				return 'bg-gray-500';
		}
	}

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

	async function handleRefresh() {
		isRefreshing = true;
		await new Promise((resolve) => setTimeout(resolve, 1000));
		isRefreshing = false;
	}
</script>

<div class="flex-1 overflow-y-auto">
	<div class="border-b border-border bg-background sticky top-0 z-10">
		<div class="flex items-center justify-between px-8 py-6">
			<div>
				<h1 class="text-3xl font-bold">Observability</h1>
				<p class="text-muted-foreground mt-1">Monitor system performance and health</p>
			</div>
			<div class="flex items-center gap-3">
				<Select.Root
					selected={{ value: selectedEnvironment, label: selectedEnvironment }}
					onSelectedChange={(v) => v && (selectedEnvironment = v.value)}
				>
					<Select.Trigger class="w-[180px]">
						<Select.Value placeholder="Select environment" />
					</Select.Trigger>
					<Select.Content>
						{#each environments as env}
							<Select.Item value={env.value}>{env.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>

				<Select.Root
					selected={{ value: timeRange, label: timeRange }}
					onSelectedChange={(v) => v && (timeRange = v.value)}
				>
					<Select.Trigger class="w-[180px]">
						<Calendar class="mr-2 h-4 w-4" />
						<Select.Value placeholder="Select time range" />
					</Select.Trigger>
					<Select.Content>
						{#each timeRanges as range}
							<Select.Item value={range.value}>{range.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>

			<Button variant="outline" size="icon" onclick={handleRefresh} disabled={isRefreshing}>
				<RefreshCw class={`h-4 w-4 ${isRefreshing ? 'animate-spin' : ''}`} />
			</Button>

				<Button variant="ghost" size="icon">
					<MoreVertical class="h-4 w-4" />
				</Button>
			</div>
		</div>
	</div>

	<div class="p-8 space-y-8">
		<section>
			<h2 class="text-xl font-semibold mb-4">Container Metrics</h2>
			<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-2">
				{#each dummyContainerMetrics as container}
					<Card>
						<CardHeader>
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-2">
									<div class={`h-2 w-2 rounded-full ${getStateColor(container.state)}`}></div>
									<CardTitle class="text-base">{container.name}</CardTitle>
								</div>
								<Badge variant="secondary">{container.state}</Badge>
							</div>
							<CardDescription>Uptime: {container.uptime}</CardDescription>
						</CardHeader>
						<CardContent class="space-y-3">
							<div class="space-y-1">
								<div class="flex justify-between text-sm">
									<span class="text-muted-foreground">CPU Usage</span>
									<span class="font-medium">{container.cpu}%</span>
								</div>
								<div class="h-2 bg-secondary rounded-full overflow-hidden">
									<div class="h-full bg-primary" style="width: {container.cpu}%"></div>
								</div>
							</div>

							<div class="space-y-1">
								<div class="flex justify-between text-sm">
									<span class="text-muted-foreground">Memory</span>
									<span class="font-medium"
										>{container.memory.used} / {container.memory.limit}</span
									>
								</div>
								<div class="h-2 bg-secondary rounded-full overflow-hidden">
									<div
										class="h-full bg-blue-500"
										style="width: {container.memory.percent}%"
									></div>
								</div>
							</div>

							<div class="grid grid-cols-2 gap-3 pt-2 border-t">
								<div>
									<div class="text-xs text-muted-foreground">Network I/O</div>
									<div class="text-sm font-medium">
										↓ {container.network.in}<br />↑ {container.network.out}
									</div>
								</div>
								<div>
									<div class="text-xs text-muted-foreground">Disk I/O</div>
									<div class="text-sm font-medium">
										R: {container.disk.read}<br />W: {container.disk.write}
									</div>
								</div>
							</div>

							{#if container.restarts > 0}
								<div class="pt-2 border-t">
									<Badge variant="destructive">Restarts: {container.restarts}</Badge>
								</div>
							{/if}
						</CardContent>
					</Card>
				{/each}
			</div>
		</section>

		<section>
			<h2 class="text-xl font-semibold mb-4">Server Metrics</h2>
			<Card>
				<CardContent class="pt-6">
					<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
						<div>
							<div class="text-sm text-muted-foreground mb-2">CPU Load Average</div>
							<div class="text-2xl font-bold">
								{serverMetrics.cpuLoadAvg[0].toFixed(2)}
							</div>
							<div class="text-xs text-muted-foreground mt-1">
								1m: {serverMetrics.cpuLoadAvg[0]} | 5m: {serverMetrics.cpuLoadAvg[1]} | 15m: {serverMetrics
									.cpuLoadAvg[2]}
							</div>
						</div>

						<div>
							<div class="text-sm text-muted-foreground mb-2">Memory</div>
							<div class="text-2xl font-bold">{serverMetrics.memory.percent}%</div>
							<div class="text-xs text-muted-foreground mt-1">
								{serverMetrics.memory.used} / {serverMetrics.memory.total}
							</div>
							<div class="h-2 bg-secondary rounded-full overflow-hidden mt-2">
								<div class="h-full bg-primary" style="width: {serverMetrics.memory.percent}%"></div>
							</div>
						</div>

						<div>
							<div class="text-sm text-muted-foreground mb-2">Disk Space</div>
							<div class="text-2xl font-bold">{serverMetrics.disk.percent}%</div>
							<div class="text-xs text-muted-foreground mt-1">
								{serverMetrics.disk.used} / {serverMetrics.disk.total}
							</div>
							<div class="h-2 bg-secondary rounded-full overflow-hidden mt-2">
								<div class="h-full bg-primary" style="width: {serverMetrics.disk.percent}%"></div>
							</div>
						</div>

						<div>
							<div class="text-sm text-muted-foreground mb-2">Uptime</div>
							<div class="text-2xl font-bold">{serverMetrics.uptime.split(' ')[0]}</div>
							<div class="text-xs text-muted-foreground mt-1">{serverMetrics.uptime}</div>
						</div>
					</div>
				</CardContent>
			</Card>
		</section>

		<section>
			<h2 class="text-xl font-semibold mb-4">Traefik Proxy Metrics</h2>
			<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
				<Card>
					<CardHeader>
						<CardDescription>Total Requests</CardDescription>
						<CardTitle class="text-3xl">{traefikMetrics.totalRequests.toLocaleString()}</CardTitle>
					</CardHeader>
					<CardContent>
						<div class="text-sm text-muted-foreground">
							{traefikMetrics.requestsPerSec} req/s
						</div>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardDescription>Average Latency</CardDescription>
						<CardTitle class="text-3xl">{traefikMetrics.avgLatency}ms</CardTitle>
					</CardHeader>
					<CardContent>
						<div class="text-sm text-muted-foreground">
							P95: {traefikMetrics.p95Latency}ms | P99: {traefikMetrics.p99Latency}ms
						</div>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardDescription>Error Rate</CardDescription>
						<CardTitle class="text-3xl">{traefikMetrics.errorRate}%</CardTitle>
					</CardHeader>
					<CardContent>
						<div class="text-sm text-muted-foreground">Last 24 hours</div>
					</CardContent>
				</Card>
			</div>
		</section>

		<section>
			<h2 class="text-xl font-semibold mb-4">PaaS Service Health</h2>
			<Card>
				<CardContent class="p-0">
					<div class="divide-y">
						{#each paasServices as service}
							<div class="p-4 hover:bg-accent/50 transition-colors">
								<div class="flex items-center justify-between">
									<div class="flex items-center gap-3">
										<div
											class="h-3 w-3 rounded-full {service.status === 'healthy'
												? 'bg-green-500'
												: service.status === 'degraded'
													? 'bg-yellow-500'
													: 'bg-red-500'}"
										></div>
										<div>
											<div class="font-medium">{service.name}</div>
											<div class="text-sm text-muted-foreground">
												{service.requests.toLocaleString()} requests
											</div>
										</div>
									</div>
									<div class="flex items-center gap-6 text-sm">
										<div>
											<div class="text-muted-foreground">Latency</div>
											<div class="font-medium">{service.latency}ms</div>
										</div>
										<div>
											<div class="text-muted-foreground">Errors</div>
											<div class="font-medium {service.errors > 10 ? 'text-red-500' : ''}">
												{service.errors}
											</div>
										</div>
										<Badge variant={service.status === 'healthy' ? 'default' : 'destructive'}>
											{service.status}
										</Badge>
									</div>
								</div>
							</div>
						{/each}
					</div>
				</CardContent>
			</Card>
		</section>

		<section>
			<div class="flex items-center justify-between mb-4">
				<h2 class="text-xl font-semibold">Quick Links</h2>
			</div>
			<div class="grid gap-4 md:grid-cols-4">
				<Card class="hover:bg-accent/50 transition-colors cursor-pointer">
					<a href="/dashboard/monitoring/deployments" class="block p-6">
						<CardTitle class="text-base">Deployment History</CardTitle>
						<CardDescription class="mt-2">View recent deployments and timings</CardDescription>
					</a>
				</Card>

				<Card class="hover:bg-accent/50 transition-colors cursor-pointer">
					<a href="/dashboard/monitoring/uptime" class="block p-6">
						<CardTitle class="text-base">Uptime Checks</CardTitle>
						<CardDescription class="mt-2">Endpoint monitoring and alerts</CardDescription>
					</a>
				</Card>

				<Card class="hover:bg-accent/50 transition-colors cursor-pointer">
					<a href="/dashboard/monitoring/logs" class="block p-6">
						<CardTitle class="text-base">Logs</CardTitle>
						<CardDescription class="mt-2">Search and filter system logs</CardDescription>
					</a>
				</Card>

				<Card class="hover:bg-accent/50 transition-colors cursor-pointer">
					<a href="/dashboard/monitoring/queues" class="block p-6">
						<CardTitle class="text-base">Queue Metrics</CardTitle>
						<CardDescription class="mt-2">Job status and performance</CardDescription>
					</a>
				</Card>
			</div>
		</section>
	</div>
</div>
