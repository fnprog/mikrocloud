<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Plus, Ellipsis, Server, Cpu, HardDrive, Activity, Loader2 } from 'lucide-svelte';
	import { createServersListQuery } from '$lib/features/servers/queries';

	const serversQuery = createServersListQuery();

	function getStatusColor(
		status: 'online' | 'offline' | 'maintenance' | 'error' | 'unknown'
	): string {
		switch (status) {
			case 'online':
				return 'bg-success';
			case 'offline':
				return 'bg-destructive';
			case 'maintenance':
				return 'bg-warning';
			case 'error':
				return 'bg-destructive';
			default:
				return 'bg-muted';
		}
	}
</script>

<svelte:head>
	<title>Servers - Dashboard</title>
</svelte:head>

<div class="container mx-auto space-y-6 p-6">
	<!-- Main Content -->
	<div class="flex-1 flex flex-col overflow-hidden">
		<!-- Header -->
		<div class="px-6 py-4">
			<div class="flex items-center justify-between">
				<div>
					<h1 class="text-2xl font-semibold">Servers</h1>
					<p class="text-sm text-muted-foreground mt-1">
						Monitor and manage your server infrastructure.
					</p>
				</div>
				<Button>
					<Plus class="w-4 h-4 mr-2" />
					New Server
				</Button>
			</div>
		</div>

		<!-- Servers Grid -->
		<div class="flex-1 overflow-auto p-6">
			{#if serversQuery.isLoading}
				<div class="flex items-center justify-center h-64">
					<Loader2 class="w-8 h-8 animate-spin text-muted-foreground" />
				</div>
			{:else if serversQuery.isError}
				<div class="flex flex-col items-center justify-center h-64 space-y-4">
					<p class="text-destructive">Failed to load servers</p>
					<Button variant="outline" onclick={() => serversQuery.refetch()}>Retry</Button>
				</div>
			{:else if serversQuery.data?.length === 0}
				<div class="flex flex-col items-center justify-center h-64 space-y-4">
					<Server class="w-12 h-12 text-muted-foreground" />
					<div class="text-center">
						<h3 class="text-lg font-semibold">No servers yet</h3>
						<p class="text-sm text-muted-foreground mt-1">
							Get started by adding your first server
						</p>
					</div>
					<Button>
						<Plus class="w-4 h-4 mr-2" />
						Add Server
					</Button>
				</div>
			{:else}
				<div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
					{#each serversQuery.data ?? [] as server (server.id)}
						<Card class="hover:shadow-md transition-shadow">
							<CardHeader class="pb-3">
								<div class="flex items-center justify-between">
									<div class="flex items-center space-x-2">
										<Server class="w-5 h-5 " />
										<CardTitle class="text-lg">{server.name}</CardTitle>
									</div>
									<div class="flex items-center space-x-2">
										<div class="w-2 h-2 rounded-full {getStatusColor(server.status)}"></div>
										<Badge variant="outline" class="text-xs">
											{server.status}
										</Badge>
									</div>
								</div>
								<div class="flex items-center space-x-4 text-sm text-muted-foreground">
									<span>{server.hostname}</span>
									<span>•</span>
									<span>{server.server_type.replace('_', ' ')}</span>
									<span>•</span>
									<span>{server.ip_address}</span>
								</div>
								{#if server.description}
									<p class="text-xs text-muted-foreground mt-1">{server.description}</p>
								{/if}
							</CardHeader>
							<CardContent>
								<div class="space-y-4">
									<!-- Server Info -->
									<div class="space-y-2 text-sm">
										{#if server.cpu_cores}
											<div class="flex items-center justify-between">
												<div class="flex items-center space-x-1">
													<Cpu class="w-4 h-4 text-stone-500" />
													<span class="text-muted-foreground">CPU Cores</span>
												</div>
												<span class="font-medium">{server.cpu_cores}</span>
											</div>
										{/if}

										{#if server.memory_mb}
											<div class="flex items-center justify-between">
												<div class="flex items-center space-x-1">
													<Activity class="w-4 h-4 text-stone-500" />
													<span class="text-muted-foreground">Memory</span>
												</div>
												<span class="font-medium">{Math.round(server.memory_mb / 1024)} GB</span>
											</div>
										{/if}

										{#if server.disk_gb}
											<div class="flex items-center justify-between">
												<div class="flex items-center space-x-1">
													<HardDrive class="w-4 h-4 text-stone-500" />
													<span class="text-muted-foreground">Disk</span>
												</div>
												<span class="font-medium">{server.disk_gb} GB</span>
											</div>
										{/if}

										{#if server.os}
											<div class="flex items-center justify-between">
												<span class="text-muted-foreground">OS</span>
												<span class="font-medium">{server.os} {server.os_version || ''}</span>
											</div>
										{/if}
									</div>

									<!-- Tags -->
									{#if server.tags && server.tags.length > 0}
										<div class="flex flex-wrap gap-1">
											{#each server.tags as tag}
												<Badge variant="secondary" class="text-xs">{tag}</Badge>
											{/each}
										</div>
									{/if}

									<!-- Actions -->
									<div class="flex space-x-2">
										<Button size="sm" variant="outline" class="flex-1">Connect</Button>
										<Button size="sm" variant="outline">
											<Ellipsis class="w-4 h-4" />
										</Button>
									</div>
								</div>
							</CardContent>
						</Card>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>
