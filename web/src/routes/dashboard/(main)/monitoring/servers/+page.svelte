<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import { Search, RefreshCw, Server, Cpu, HardDrive, Network, Activity, AlertTriangle, Clock, CheckCircle, XCircle, ChevronDown, ChevronUp } from 'lucide-svelte';

	let searchQuery = $state('');
	let statusFilter = $state('all');
	let regionFilter = $state('all');
	let isRefreshing = $state(false);
	let expandedServer = $state<string | null>(null);

	const statusFilters = [
		{ value: 'all', label: 'All Status' },
		{ value: 'online', label: 'Online' },
		{ value: 'offline', label: 'Offline' },
		{ value: 'warning', label: 'Warning' }
	];

	const regionFilters = [
		{ value: 'all', label: 'All Regions' },
		{ value: 'us-east-1', label: 'US East' },
		{ value: 'us-west-2', label: 'US West' },
		{ value: 'eu-central-1', label: 'EU Central' },
		{ value: 'ap-south-1', label: 'AP South' }
	];

	const servers = [
		{
			id: 'srv_001',
			hostname: 'web-prod-01',
			ip: '10.0.1.15',
			region: 'us-east-1',
			status: 'online',
			uptime: 2592000,
			cpu: { usage: 42, cores: 8, model: 'Intel Xeon E5-2686 v4' },
			memory: { usage: 68, total: 32, used: 21.76 },
			disk: { usage: 55, total: 500, used: 275 },
			network: { in: 1250, out: 3420 },
			containers: 12,
			lastSeen: '2024-10-10T14:35:00Z',
			tags: ['production', 'web', 'high-priority']
		},
		{
			id: 'srv_002',
			hostname: 'web-prod-02',
			ip: '10.0.1.16',
			region: 'us-east-1',
			status: 'online',
			uptime: 1987200,
			cpu: { usage: 38, cores: 8, model: 'Intel Xeon E5-2686 v4' },
			memory: { usage: 62, total: 32, used: 19.84 },
			disk: { usage: 48, total: 500, used: 240 },
			network: { in: 1120, out: 3150 },
			containers: 11,
			lastSeen: '2024-10-10T14:35:05Z',
			tags: ['production', 'web']
		},
		{
			id: 'srv_003',
			hostname: 'db-primary-01',
			ip: '10.0.2.10',
			region: 'us-east-1',
			status: 'warning',
			uptime: 3456000,
			cpu: { usage: 78, cores: 16, model: 'AMD EPYC 7763' },
			memory: { usage: 89, total: 64, used: 56.96 },
			disk: { usage: 82, total: 2000, used: 1640 },
			network: { in: 5240, out: 4890 },
			containers: 3,
			lastSeen: '2024-10-10T14:35:10Z',
			tags: ['production', 'database', 'critical']
		},
		{
			id: 'srv_004',
			hostname: 'api-prod-01',
			ip: '10.0.3.20',
			region: 'us-west-2',
			status: 'online',
			uptime: 1728000,
			cpu: { usage: 52, cores: 12, model: 'Intel Xeon Platinum 8375C' },
			memory: { usage: 71, total: 48, used: 34.08 },
			disk: { usage: 38, total: 1000, used: 380 },
			network: { in: 3450, out: 6720 },
			containers: 18,
			lastSeen: '2024-10-10T14:35:15Z',
			tags: ['production', 'api']
		},
		{
			id: 'srv_005',
			hostname: 'cache-prod-01',
			ip: '10.0.4.30',
			region: 'us-east-1',
			status: 'online',
			uptime: 4320000,
			cpu: { usage: 24, cores: 4, model: 'Intel Xeon E5-2686 v4' },
			memory: { usage: 95, total: 16, used: 15.2 },
			disk: { usage: 22, total: 250, used: 55 },
			network: { in: 8920, out: 7650 },
			containers: 2,
			lastSeen: '2024-10-10T14:35:20Z',
			tags: ['production', 'cache', 'redis']
		},
		{
			id: 'srv_006',
			hostname: 'worker-prod-01',
			ip: '10.0.5.40',
			region: 'eu-central-1',
			status: 'online',
			uptime: 864000,
			cpu: { usage: 66, cores: 16, model: 'AMD EPYC 7R13' },
			memory: { usage: 54, total: 64, used: 34.56 },
			disk: { usage: 45, total: 500, used: 225 },
			network: { in: 890, out: 1230 },
			containers: 24,
			lastSeen: '2024-10-10T14:35:25Z',
			tags: ['production', 'worker', 'queue']
		},
		{
			id: 'srv_007',
			hostname: 'staging-01',
			ip: '10.0.6.50',
			region: 'us-west-2',
			status: 'warning',
			uptime: 432000,
			cpu: { usage: 92, cores: 4, model: 'Intel Xeon E5-2686 v4' },
			memory: { usage: 88, total: 16, used: 14.08 },
			disk: { usage: 76, total: 250, used: 190 },
			network: { in: 420, out: 680 },
			containers: 8,
			lastSeen: '2024-10-10T14:35:30Z',
			tags: ['staging', 'testing']
		},
		{
			id: 'srv_008',
			hostname: 'build-prod-01',
			ip: '10.0.7.60',
			region: 'us-east-1',
			status: 'online',
			uptime: 1296000,
			cpu: { usage: 48, cores: 8, model: 'Intel Xeon Platinum 8375C' },
			memory: { usage: 58, total: 32, used: 18.56 },
			disk: { usage: 68, total: 1000, used: 680 },
			network: { in: 2340, out: 1890 },
			containers: 6,
			lastSeen: '2024-10-10T14:35:35Z',
			tags: ['production', 'build', 'ci-cd']
		},
		{
			id: 'srv_009',
			hostname: 'monitoring-01',
			ip: '10.0.8.70',
			region: 'us-east-1',
			status: 'online',
			uptime: 5184000,
			cpu: { usage: 32, cores: 4, model: 'Intel Xeon E5-2686 v4' },
			memory: { usage: 64, total: 16, used: 10.24 },
			disk: { usage: 58, total: 500, used: 290 },
			network: { in: 1560, out: 2340 },
			containers: 5,
			lastSeen: '2024-10-10T14:35:40Z',
			tags: ['production', 'monitoring', 'metrics']
		},
		{
			id: 'srv_010',
			hostname: 'dev-legacy-01',
			ip: '10.0.9.80',
			region: 'ap-south-1',
			status: 'offline',
			uptime: 0,
			cpu: { usage: 0, cores: 4, model: 'Intel Xeon E5-2686 v4' },
			memory: { usage: 0, total: 16, used: 0 },
			disk: { usage: 0, total: 250, used: 0 },
			network: { in: 0, out: 0 },
			containers: 0,
			lastSeen: '2024-10-10T08:15:00Z',
			tags: ['development', 'legacy', 'deprecated']
		}
	];

	const filteredServers = $derived(
		servers.filter((server) => {
			const matchesSearch =
				server.hostname.toLowerCase().includes(searchQuery.toLowerCase()) ||
				server.ip.includes(searchQuery) ||
				server.tags.some(tag => tag.toLowerCase().includes(searchQuery.toLowerCase()));
			const matchesStatus = statusFilter === 'all' || server.status === statusFilter;
			const matchesRegion = regionFilter === 'all' || server.region === regionFilter;
			return matchesSearch && matchesStatus && matchesRegion;
		})
	);

	const stats = $derived({
		total: servers.length,
		online: servers.filter((s) => s.status === 'online').length,
		offline: servers.filter((s) => s.status === 'offline').length,
		warning: servers.filter((s) => s.status === 'warning').length,
		avgCpu: Math.round(servers.filter(s => s.status !== 'offline').reduce((sum, s) => sum + s.cpu.usage, 0) / servers.filter(s => s.status !== 'offline').length),
		avgMemory: Math.round(servers.filter(s => s.status !== 'offline').reduce((sum, s) => sum + s.memory.usage, 0) / servers.filter(s => s.status !== 'offline').length),
		totalContainers: servers.reduce((sum, s) => sum + s.containers, 0)
	});

	function getStatusColor(status: string) {
		switch (status) {
			case 'online':
				return 'text-green-500';
			case 'offline':
				return 'text-red-500';
			case 'warning':
				return 'text-yellow-500';
			default:
				return 'text-gray-500';
		}
	}

	function getStatusBgColor(status: string) {
		switch (status) {
			case 'online':
				return 'bg-green-500/10 border-green-500/20';
			case 'offline':
				return 'bg-red-500/10 border-red-500/20';
			case 'warning':
				return 'bg-yellow-500/10 border-yellow-500/20';
			default:
				return 'bg-gray-500/10 border-gray-500/20';
		}
	}

	function getStatusIcon(status: string) {
		switch (status) {
			case 'online':
				return CheckCircle;
			case 'offline':
				return XCircle;
			case 'warning':
				return AlertTriangle;
			default:
				return Activity;
		}
	}

	function formatUptime(seconds: number) {
		if (seconds === 0) return 'Offline';
		const days = Math.floor(seconds / 86400);
		const hours = Math.floor((seconds % 86400) / 3600);
		if (days > 0) return `${days}d ${hours}h`;
		return `${hours}h`;
	}

	function formatBytes(gb: number) {
		return `${gb.toFixed(2)} GB`;
	}

	function formatBandwidth(kbps: number) {
		if (kbps < 1000) return `${kbps} KB/s`;
		return `${(kbps / 1000).toFixed(2)} MB/s`;
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
		expandedServer = expandedServer === id ? null : id;
	}

	async function handleRefresh() {
		isRefreshing = true;
		await new Promise((resolve) => setTimeout(resolve, 1000));
		isRefreshing = false;
	}

	function exportServers() {
		const dataStr = JSON.stringify(filteredServers, null, 2);
		const blob = new Blob([dataStr], { type: 'application/json' });
		const url = URL.createObjectURL(blob);
		const link = document.createElement('a');
		link.href = url;
		link.download = `servers-${new Date().toISOString()}.json`;
		link.click();
		URL.revokeObjectURL(url);
	}
</script>

<div class="flex-1 overflow-y-auto">
	<div class="border-b border-border bg-background sticky top-0 z-10">
		<div class="px-8 py-6">
			<div class="flex items-center justify-between mb-4">
				<div>
					<h1 class="text-3xl font-bold">Servers</h1>
					<p class="text-muted-foreground mt-1">Monitor server health and resource usage</p>
				</div>
				<div class="flex items-center gap-2">
					<Button variant="outline" size="sm" onclick={exportServers}>
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
						<CardDescription class="text-xs">Total Servers</CardDescription>
						<CardTitle class="text-2xl">{stats.total}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Online</CardDescription>
						<CardTitle class="text-2xl text-green-500">{stats.online}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Offline</CardDescription>
						<CardTitle class="text-2xl text-red-500">{stats.offline}</CardTitle>
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
						<CardDescription class="text-xs">Avg CPU</CardDescription>
						<CardTitle class="text-2xl text-blue-500">{stats.avgCpu}%</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Avg Memory</CardDescription>
						<CardTitle class="text-2xl text-purple-500">{stats.avgMemory}%</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Containers</CardDescription>
						<CardTitle class="text-2xl text-orange-500">{stats.totalContainers}</CardTitle>
					</CardHeader>
				</Card>
			</div>

			<div class="grid grid-cols-3 gap-3">
				<div class="relative">
					<Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
					<Input
						type="text"
						placeholder="Search servers..."
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
					selected={{ value: regionFilter, label: regionFilter }}
					onSelectedChange={(v) => v && (regionFilter = v.value)}
				>
					<Select.Trigger>
						<Select.Value placeholder="Filter by region" />
					</Select.Trigger>
					<Select.Content>
						{#each regionFilters as filter}
							<Select.Item value={filter.value}>{filter.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		</div>
	</div>

	<div class="p-8 space-y-3">
		<div class="text-sm text-muted-foreground mb-2">
			Showing {filteredServers.length} of {servers.length} servers
		</div>

	{#each filteredServers as server (server.id)}
		{@const Icon = getStatusIcon(server.status)}
		<Card class="hover:bg-muted/50 transition-colors">
			<CardContent class="p-4">
				<div class="flex items-start gap-4">
					<div class="p-2 border rounded-lg {getStatusBgColor(server.status)}">
						<Icon class="h-5 w-5 {getStatusColor(server.status)}" />
					</div>

						<div class="flex-1 min-w-0">
							<div class="flex items-start justify-between gap-4 mb-3">
								<div class="flex-1">
									<div class="flex items-center gap-2 mb-1">
										<span class="font-semibold text-lg">{server.hostname}</span>
										<Badge variant="outline" class="text-xs font-mono">{server.ip}</Badge>
										<Badge variant={server.status === 'online' ? 'default' : server.status === 'offline' ? 'destructive' : 'secondary'} class="capitalize">
											{server.status}
										</Badge>
									</div>
									<div class="flex items-center gap-3 text-sm text-muted-foreground mb-2">
										<span>Region: <span class="font-mono">{server.region}</span></span>
										<span>•</span>
										<span>Uptime: {formatUptime(server.uptime)}</span>
										<span>•</span>
										<span>Containers: {server.containers}</span>
									</div>
									<div class="flex gap-1">
										{#each server.tags as tag}
											<Badge variant="outline" class="text-xs">{tag}</Badge>
										{/each}
									</div>
								</div>

								<Button
									variant="ghost"
									size="sm"
									onclick={() => toggleExpanded(server.id)}
								>
									{#if expandedServer === server.id}
										<ChevronUp class="h-4 w-4" />
									{:else}
										<ChevronDown class="h-4 w-4" />
									{/if}
								</Button>
							</div>

							<div class="grid grid-cols-4 gap-4">
								<div>
									<div class="flex items-center gap-2 mb-2">
										<Cpu class="h-4 w-4 text-muted-foreground" />
										<span class="text-xs text-muted-foreground">CPU Usage</span>
									</div>
									<div class="space-y-1">
										<div class="flex items-center justify-between text-sm">
											<span class="font-semibold {getResourceColor(server.cpu.usage)}">{server.cpu.usage}%</span>
											<span class="text-xs text-muted-foreground">{server.cpu.cores} cores</span>
										</div>
										<div class="h-2 bg-muted rounded-full overflow-hidden">
											<div class="{getResourceBgColor(server.cpu.usage)} h-full transition-all" style="width: {server.cpu.usage}%"></div>
										</div>
									</div>
								</div>

								<div>
									<div class="flex items-center gap-2 mb-2">
										<Activity class="h-4 w-4 text-muted-foreground" />
										<span class="text-xs text-muted-foreground">Memory Usage</span>
									</div>
									<div class="space-y-1">
										<div class="flex items-center justify-between text-sm">
											<span class="font-semibold {getResourceColor(server.memory.usage)}">{server.memory.usage}%</span>
											<span class="text-xs text-muted-foreground">{formatBytes(server.memory.used)} / {formatBytes(server.memory.total)}</span>
										</div>
										<div class="h-2 bg-muted rounded-full overflow-hidden">
											<div class="{getResourceBgColor(server.memory.usage)} h-full transition-all" style="width: {server.memory.usage}%"></div>
										</div>
									</div>
								</div>

								<div>
									<div class="flex items-center gap-2 mb-2">
										<HardDrive class="h-4 w-4 text-muted-foreground" />
										<span class="text-xs text-muted-foreground">Disk Usage</span>
									</div>
									<div class="space-y-1">
										<div class="flex items-center justify-between text-sm">
											<span class="font-semibold {getResourceColor(server.disk.usage)}">{server.disk.usage}%</span>
											<span class="text-xs text-muted-foreground">{server.disk.used} / {server.disk.total} GB</span>
										</div>
										<div class="h-2 bg-muted rounded-full overflow-hidden">
											<div class="{getResourceBgColor(server.disk.usage)} h-full transition-all" style="width: {server.disk.usage}%"></div>
										</div>
									</div>
								</div>

								<div>
									<div class="flex items-center gap-2 mb-2">
										<Network class="h-4 w-4 text-muted-foreground" />
										<span class="text-xs text-muted-foreground">Network I/O</span>
									</div>
									<div class="space-y-1">
										<div class="flex items-center justify-between text-sm">
											<div class="flex items-center gap-1">
												<span class="text-xs text-muted-foreground">↓</span>
												<span class="font-semibold text-green-500">{formatBandwidth(server.network.in)}</span>
											</div>
											<div class="flex items-center gap-1">
												<span class="text-xs text-muted-foreground">↑</span>
												<span class="font-semibold text-blue-500">{formatBandwidth(server.network.out)}</span>
											</div>
										</div>
									</div>
								</div>
							</div>

							<div class="text-xs text-muted-foreground mt-3">
								Last seen: {formatTimestamp(server.lastSeen)} ({getRelativeTime(server.lastSeen)})
							</div>

							{#if expandedServer === server.id}
								<div class="mt-4 pt-4 border-t">
									<div class="text-sm font-semibold mb-3">Server Details</div>
									<div class="grid grid-cols-2 gap-4">
										<div class="space-y-2">
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">Hostname:</span>
												<span class="font-mono">{server.hostname}</span>
											</div>
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">IP Address:</span>
												<span class="font-mono">{server.ip}</span>
											</div>
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">Region:</span>
												<span class="font-mono">{server.region}</span>
											</div>
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">Status:</span>
												<Badge variant={server.status === 'online' ? 'default' : server.status === 'offline' ? 'destructive' : 'secondary'} class="capitalize">
													{server.status}
												</Badge>
											</div>
										</div>
										<div class="space-y-2">
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">CPU Model:</span>
												<span class="text-xs">{server.cpu.model}</span>
											</div>
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">CPU Cores:</span>
												<span>{server.cpu.cores}</span>
											</div>
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">Total Memory:</span>
												<span>{formatBytes(server.memory.total)}</span>
											</div>
											<div class="flex justify-between text-sm">
												<span class="text-muted-foreground">Total Disk:</span>
												<span>{server.disk.total} GB</span>
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

		{#if filteredServers.length === 0}
			<Card>
				<CardContent class="py-12">
					<div class="text-center text-muted-foreground">
						No servers found matching your criteria
					</div>
				</CardContent>
			</Card>
		{/if}
	</div>
</div>
