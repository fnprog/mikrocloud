<script lang="ts">
	import { createQuery } from '@tanstack/svelte-query';
	import { serversApi, type Server } from '$lib/api/servers';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Server as ServerIcon, Plus, Activity, Cpu, HardDrive, LoaderCircle } from 'lucide-svelte';

	const serversQuery = createQuery(() => ({
		queryKey: ['servers'],
		queryFn: () => serversApi.list(),
		staleTime: 30 * 1000
	}));

	function getStatusColor(status: Server['status']) {
		switch (status) {
			case 'online':
				return 'bg-green-500';
			case 'offline':
				return 'bg-gray-500';
			case 'maintenance':
				return 'bg-yellow-500';
			case 'error':
				return 'bg-red-500';
			default:
				return 'bg-gray-400';
		}
	}

	function getServerTypeLabel(type: Server['server_type']) {
		return type
			.split('_')
			.map((word) => word.charAt(0).toUpperCase() + word.slice(1))
			.join(' ');
	}

	function formatBytes(mb: number | undefined) {
		if (!mb) return 'N/A';
		if (mb >= 1024) return `${(mb / 1024).toFixed(1)} GB`;
		return `${mb} MB`;
	}
</script>

<svelte:head>
	<title>Infrastructure - Dashboard</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold">Infrastructure</h1>
			<p class="text-muted-foreground mt-1">Manage your servers and infrastructure resources</p>
		</div>
		<Button>
			<Plus class="h-4 w-4 mr-2" />
			Add Server
		</Button>
	</div>

	{#if serversQuery.isLoading}
		<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
			{#each [1, 2, 3] as i}
				<Card>
					<CardHeader>
						<div class="h-4 bg-muted rounded animate-pulse w-1/2"></div>
						<div class="h-3 bg-muted rounded animate-pulse w-1/3 mt-2"></div>
					</CardHeader>
					<CardContent>
						<div class="h-3 bg-muted rounded animate-pulse w-full"></div>
					</CardContent>
				</Card>
			{/each}
		</div>
	{:else if serversQuery.error}
		<Card class="border-destructive/50">
			<CardContent class="pt-6">
				<p class="text-destructive">
					Error loading servers: {serversQuery.error instanceof Error ? serversQuery.error.message : 'Unknown error'}
				</p>
			</CardContent>
		</Card>
	{:else if serversQuery.data && serversQuery.data.length === 0}
		<Card>
			<CardContent class="pt-6 text-center py-12">
				<ServerIcon class="h-12 w-12 mx-auto text-muted-foreground mb-4" />
				<h3 class="text-lg font-medium mb-2">No servers yet</h3>
				<p class="text-muted-foreground mb-4">
					Add your first server to start deploying applications
				</p>
				<Button>
					<Plus class="h-4 w-4 mr-2" />
					Add Your First Server
				</Button>
			</CardContent>
		</Card>
	{:else if serversQuery.data}
		<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
			{#each serversQuery.data as server (server.id)}
				<Card class="hover:shadow-md transition-shadow">
					<CardHeader>
						<div class="flex items-start justify-between">
							<div class="flex items-center gap-3 flex-1 min-w-0">
								<div class="p-2 rounded-lg bg-primary/10">
									<ServerIcon class="h-5 w-5 text-primary" />
								</div>
								<div class="flex-1 min-w-0">
									<CardTitle class="text-lg truncate">{server.name}</CardTitle>
									<div class="flex items-center gap-2 mt-1">
										<div class="flex items-center gap-1.5">
											<div class="w-2 h-2 rounded-full {getStatusColor(server.status)}"></div>
											<span class="text-xs text-muted-foreground capitalize">{server.status}</span>
										</div>
										<Badge variant="secondary" class="text-xs">
											{getServerTypeLabel(server.server_type)}
										</Badge>
									</div>
								</div>
							</div>
						</div>
					</CardHeader>
					<CardContent>
						{#if server.description}
							<p class="text-sm text-muted-foreground mb-3 line-clamp-2">{server.description}</p>
						{/if}
						<div class="space-y-2 text-sm">
							<div class="flex items-center justify-between">
								<span class="text-muted-foreground">IP Address</span>
								<span class="font-mono text-xs">{server.ip_address}</span>
							</div>
							{#if server.cpu_cores || server.memory_mb || server.disk_gb}
								<div class="flex items-center gap-4 pt-2 border-t">
									{#if server.cpu_cores}
										<div class="flex items-center gap-1.5">
											<Cpu class="h-3.5 w-3.5 text-muted-foreground" />
											<span class="text-xs">{server.cpu_cores} cores</span>
										</div>
									{/if}
									{#if server.memory_mb}
										<div class="flex items-center gap-1.5">
											<Activity class="h-3.5 w-3.5 text-muted-foreground" />
											<span class="text-xs">{formatBytes(server.memory_mb)}</span>
										</div>
									{/if}
									{#if server.disk_gb}
										<div class="flex items-center gap-1.5">
											<HardDrive class="h-3.5 w-3.5 text-muted-foreground" />
											<span class="text-xs">{server.disk_gb} GB</span>
										</div>
									{/if}
								</div>
							{/if}
							{#if server.tags && server.tags.length > 0}
								<div class="flex items-center gap-1.5 flex-wrap pt-2">
									{#each server.tags as tag}
										<Badge variant="outline" class="text-xs">{tag}</Badge>
									{/each}
								</div>
							{/if}
						</div>
					</CardContent>
				</Card>
			{/each}
		</div>
	{/if}
</div>
