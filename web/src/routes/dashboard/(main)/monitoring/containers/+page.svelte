<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import { Search, RefreshCw, Play, Square, Trash2, Terminal, MoreVertical } from 'lucide-svelte';

	let searchQuery = $state('');
	let statusFilter = $state('all');
	let isRefreshing = $state(false);

	const containers = [
		{
			id: 'mikrocloud-api',
			name: 'mikrocloud-api',
			image: 'mikrocloud/api:latest',
			status: 'running',
			state: 'Up 15 days',
			cpu: 45.2,
			cpuHistory: [42, 43, 45, 46, 45, 47, 45, 44, 45, 46],
			memory: { used: 512, limit: 1024, percent: 50 },
			memoryHistory: [48, 49, 50, 51, 50, 50, 49, 50, 51, 50],
			network: { rx: 1.2, tx: 0.8, rxTotal: '1.5GB', txTotal: '980MB' },
			disk: { read: 50, write: 25, readTotal: '12GB', writeTotal: '8GB' },
			ports: ['80:8080', '443:8443'],
			restarts: 0,
			created: '2024-09-25T10:30:00Z'
		},
		{
			id: 'mikrocloud-scheduler',
			name: 'mikrocloud-scheduler',
			image: 'mikrocloud/scheduler:v1.2',
			status: 'running',
			state: 'Up 15 days',
			cpu: 12.5,
			cpuHistory: [11, 12, 13, 12, 12, 13, 12, 11, 12, 13],
			memory: { used: 256, limit: 512, percent: 50 },
			memoryHistory: [49, 50, 51, 50, 49, 50, 51, 50, 49, 50],
			network: { rx: 0.1, tx: 0.05, rxTotal: '250MB', txTotal: '120MB' },
			disk: { read: 10, write: 5, readTotal: '2GB', writeTotal: '1GB' },
			ports: ['8081:8081'],
			restarts: 1,
			created: '2024-09-25T10:28:00Z'
		},
		{
			id: 'app-web-1',
			name: 'app-web-1',
			image: 'myapp/web:prod-abc123',
			status: 'running',
			state: 'Up 7 days',
			cpu: 78.9,
			cpuHistory: [75, 76, 78, 80, 79, 78, 77, 78, 79, 79],
			memory: { used: 1536, limit: 2048, percent: 75 },
			memoryHistory: [73, 74, 75, 76, 75, 74, 75, 76, 75, 75],
			network: { rx: 5.5, tx: 3.2, rxTotal: '45GB', txTotal: '28GB' },
			disk: { read: 200, write: 150, readTotal: '80GB', writeTotal: '55GB' },
			ports: ['8082:3000'],
			restarts: 0,
			created: '2024-10-03T08:15:00Z'
		},
		{
			id: 'postgres-main',
			name: 'postgres-main',
			image: 'postgres:16-alpine',
			status: 'running',
			state: 'Up 30 days',
			cpu: 32.1,
			cpuHistory: [30, 31, 32, 33, 32, 31, 32, 33, 32, 32],
			memory: { used: 768, limit: 2048, percent: 38 },
			memoryHistory: [37, 38, 39, 38, 37, 38, 39, 38, 37, 38],
			network: { rx: 2.1, tx: 1.8, rxTotal: '120GB', txTotal: '98GB' },
			disk: { read: 500, write: 350, readTotal: '500GB', writeTotal: '380GB' },
			ports: ['5432:5432'],
			restarts: 0,
			created: '2024-09-10T12:00:00Z'
		},
		{
			id: 'redis-cache',
			name: 'redis-cache',
			image: 'redis:7-alpine',
			status: 'running',
			state: 'Up 20 days',
			cpu: 8.3,
			cpuHistory: [7, 8, 9, 8, 7, 8, 9, 8, 8, 8],
			memory: { used: 128, limit: 512, percent: 25 },
			memoryHistory: [24, 25, 26, 25, 24, 25, 26, 25, 24, 25],
			network: { rx: 0.5, tx: 0.4, rxTotal: '8GB', txTotal: '6GB' },
			disk: { read: 5, write: 3, readTotal: '500MB', writeTotal: '300MB' },
			ports: ['6379:6379'],
			restarts: 0,
			created: '2024-09-20T15:20:00Z'
		},
		{
			id: 'nginx-proxy',
			name: 'nginx-proxy',
			image: 'nginx:alpine',
			status: 'exited',
			state: 'Exited (1) 2 hours ago',
			cpu: 0,
			cpuHistory: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
			memory: { used: 0, limit: 256, percent: 0 },
			memoryHistory: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
			network: { rx: 0, tx: 0, rxTotal: '1.2GB', txTotal: '800MB' },
			disk: { read: 0, write: 0, readTotal: '50MB', writeTotal: '30MB' },
			ports: ['80:80', '443:443'],
			restarts: 3,
			created: '2024-10-08T14:00:00Z'
		}
	];

	const filteredContainers = $derived(
		containers.filter((c) => {
			const matchesSearch = c.name.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesStatus = statusFilter === 'all' || c.status === statusFilter;
			return matchesSearch && matchesStatus;
		})
	);

	function getStatusColor(status: string) {
		switch (status) {
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

	function getCpuColor(cpu: number) {
		if (cpu < 50) return 'bg-green-500';
		if (cpu < 80) return 'bg-yellow-500';
		return 'bg-red-500';
	}

	function getMemoryColor(percent: number) {
		if (percent < 60) return 'bg-blue-500';
		if (percent < 85) return 'bg-yellow-500';
		return 'bg-red-500';
	}

	async function handleRefresh() {
		isRefreshing = true;
		await new Promise((resolve) => setTimeout(resolve, 1000));
		isRefreshing = false;
	}
</script>

<div class="flex-1 overflow-y-auto">
	<div class="border-b border-border bg-background sticky top-0 z-10">
		<div class="px-8 py-6">
			<div class="flex items-center justify-between mb-4">
				<div>
					<h1 class="text-3xl font-bold">Containers</h1>
					<p class="text-muted-foreground mt-1">Monitor and manage container resources</p>
				</div>
			<Button variant="outline" size="icon" onclick={handleRefresh} disabled={isRefreshing}>
				<RefreshCw class={`h-4 w-4 ${isRefreshing ? 'animate-spin' : ''}`} />
			</Button>
			</div>

			<div class="flex items-center gap-3">
				<div class="relative flex-1 max-w-md">
					<Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
					<Input
						type="text"
						placeholder="Search containers..."
						class="pl-10"
						bind:value={searchQuery}
					/>
				</div>

				<Select.Root
					selected={{ value: statusFilter, label: statusFilter }}
					onSelectedChange={(v) => v && (statusFilter = v.value)}
				>
					<Select.Trigger class="w-[180px]">
						<Select.Value placeholder="Filter by status" />
					</Select.Trigger>
					<Select.Content>
						<Select.Item value="all">All Statuses</Select.Item>
						<Select.Item value="running">Running</Select.Item>
						<Select.Item value="exited">Exited</Select.Item>
						<Select.Item value="paused">Paused</Select.Item>
					</Select.Content>
				</Select.Root>
			</div>
		</div>
	</div>

	<div class="p-8 space-y-4">
		<div class="flex items-center justify-between mb-2">
			<div class="text-sm text-muted-foreground">
				Showing {filteredContainers.length} of {containers.length} containers
			</div>
		</div>

		{#each filteredContainers as container (container.id)}
			<Card>
				<CardHeader>
					<div class="flex items-center justify-between">
						<div class="flex items-center gap-3">
							<div class={`h-3 w-3 rounded-full ${getStatusColor(container.status)}`}></div>
							<div>
								<CardTitle class="text-lg">{container.name}</CardTitle>
								<CardDescription class="mt-1">{container.image}</CardDescription>
							</div>
						</div>
						<div class="flex items-center gap-2">
							<Badge variant={container.status === 'running' ? 'default' : 'secondary'}>
								{container.status}
							</Badge>
							{#if container.restarts > 0}
								<Badge variant="destructive">Restarts: {container.restarts}</Badge>
							{/if}
							<Button variant="ghost" size="icon">
								<MoreVertical class="h-4 w-4" />
							</Button>
						</div>
					</div>
				</CardHeader>
				<CardContent>
					<div class="grid gap-6 lg:grid-cols-2">
						<div class="space-y-4">
							<div>
								<div class="flex items-center justify-between mb-2">
									<span class="text-sm font-medium">CPU Usage</span>
									<span class="text-sm font-bold">{container.cpu}%</span>
								</div>
								<div class="h-2 bg-secondary rounded-full overflow-hidden">
									<div
										class="h-full {getCpuColor(container.cpu)}"
										style="width: {container.cpu}%"
									></div>
								</div>
								<div class="mt-2 h-12 flex items-end gap-1">
									{#each container.cpuHistory as value}
										<div
											class="flex-1 bg-primary/30 rounded-t"
											style="height: {value}%"
										></div>
									{/each}
								</div>
							</div>

							<div>
								<div class="flex items-center justify-between mb-2">
									<span class="text-sm font-medium">Memory Usage</span>
									<span class="text-sm font-bold"
										>{container.memory.used}MB / {container.memory.limit}MB</span
									>
								</div>
								<div class="h-2 bg-secondary rounded-full overflow-hidden">
									<div
										class="h-full {getMemoryColor(container.memory.percent)}"
										style="width: {container.memory.percent}%"
									></div>
								</div>
								<div class="mt-2 h-12 flex items-end gap-1">
									{#each container.memoryHistory as value}
										<div
											class="flex-1 bg-blue-500/30 rounded-t"
											style="height: {value}%"
										></div>
									{/each}
								</div>
							</div>
						</div>

						<div class="grid grid-cols-2 gap-4">
							<div class="space-y-3">
								<div class="border rounded-lg p-3">
									<div class="text-xs text-muted-foreground mb-1">Network I/O</div>
									<div class="text-sm font-medium">
										↓ {container.network.rx} MB/s
									</div>
									<div class="text-sm font-medium">
										↑ {container.network.tx} MB/s
									</div>
									<div class="text-xs text-muted-foreground mt-2">
										Total: {container.network.rxTotal} / {container.network.txTotal}
									</div>
								</div>

								<div class="border rounded-lg p-3">
									<div class="text-xs text-muted-foreground mb-1">Disk I/O</div>
									<div class="text-sm font-medium">
										R: {container.disk.read} KB/s
									</div>
									<div class="text-sm font-medium">
										W: {container.disk.write} KB/s
									</div>
									<div class="text-xs text-muted-foreground mt-2">
										Total: {container.disk.readTotal} / {container.disk.writeTotal}
									</div>
								</div>
							</div>

							<div class="space-y-3">
								<div class="border rounded-lg p-3">
									<div class="text-xs text-muted-foreground mb-1">Ports</div>
									{#each container.ports as port}
										<div class="text-sm font-medium font-mono">{port}</div>
									{/each}
								</div>

								<div class="border rounded-lg p-3">
									<div class="text-xs text-muted-foreground mb-1">Status</div>
									<div class="text-sm font-medium">{container.state}</div>
									<div class="text-xs text-muted-foreground mt-2">
										Created: {new Date(container.created).toLocaleDateString()}
									</div>
								</div>
							</div>
						</div>
					</div>

					{#if container.status === 'running'}
						<div class="flex items-center gap-2 mt-4 pt-4 border-t">
							<Button variant="outline" size="sm">
								<Square class="h-3 w-3 mr-2" />
								Stop
							</Button>
							<Button variant="outline" size="sm">
								<Terminal class="h-3 w-3 mr-2" />
								Logs
							</Button>
							<Button variant="outline" size="sm">
								<Terminal class="h-3 w-3 mr-2" />
								Exec
							</Button>
							<Button variant="destructive" size="sm" class="ml-auto">
								<Trash2 class="h-3 w-3 mr-2" />
								Remove
							</Button>
						</div>
					{:else}
						<div class="flex items-center gap-2 mt-4 pt-4 border-t">
							<Button variant="outline" size="sm">
								<Play class="h-3 w-3 mr-2" />
								Start
							</Button>
							<Button variant="destructive" size="sm" class="ml-auto">
								<Trash2 class="h-3 w-3 mr-2" />
								Remove
							</Button>
						</div>
					{/if}
				</CardContent>
			</Card>
		{/each}

		{#if filteredContainers.length === 0}
			<Card>
				<CardContent class="py-12">
					<div class="text-center text-muted-foreground">
						No containers found matching your criteria
					</div>
				</CardContent>
			</Card>
		{/if}
	</div>
</div>
