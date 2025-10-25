<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Card from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Badge } from '$lib/components/ui/badge';
	import { Play, Square, RotateCw, Trash2, Plus, Network } from 'lucide-svelte';
	import { createTunnelsFetchQuery } from '$lib/features/tunnels/queries';
	import {
		createTunnelMutationQuery,
		deleteTunnelMutationQuery,
		startTunnelMutationQuery,
		stopTunnelMutationQuery,
		restartTunnelMutationQuery
	} from '$lib/features/tunnels/mutations';
	import type { Tunnel, TunnelStatus, HealthStatus } from '$lib/features/tunnels/types';

	const tunnelsQuery = createTunnelsFetchQuery();

	let showCreateDialog = $state(false);
	let name = $state('');
	let tunnelToken = $state('');
	let projectId = $state('');

	const createTunnelMutation = createTunnelMutationQuery({
		onSuccess: () => {
			showCreateDialog = false;
			resetForm();
		},
		onError: (error) => {
			alert(`Failed to create tunnel: ${error.message}`);
		}
	});

	const startTunnelMutation = startTunnelMutationQuery({
		onError: (error) => alert(`Failed to start tunnel: ${error.message}`)
	});

	function resetForm() {
		name = '';
		tunnelToken = '';
		projectId = '';
	}

	function handleCreate() {
		createTunnelMutation.mutate({
			name,
			tunnel_token: tunnelToken,
			project_id: projectId || undefined
		});
	}

	function handleStart(tunnel: Tunnel) {
		startTunnelMutation.mutate(tunnel.id);
	}

	function handleStop(tunnel: Tunnel) {
		const mutation = stopTunnelMutationQuery(tunnel.id, {
			onError: (error) => alert(`Failed to stop tunnel: ${error.message}`)
		});
		mutation.mutate();
	}

	function handleRestart(tunnel: Tunnel) {
		const mutation = restartTunnelMutationQuery(tunnel.id, {
			onError: (error) => alert(`Failed to restart tunnel: ${error.message}`)
		});
		mutation.mutate();
	}

	function handleDelete(tunnel: Tunnel) {
		if (!confirm(`Are you sure you want to delete tunnel "${tunnel.name}"?`)) return;

		const mutation = deleteTunnelMutationQuery(tunnel.id, {
			onError: (error) => alert(`Failed to delete tunnel: ${error.message}`)
		});
		mutation.mutate();
	}

	function getStatusColor(status: TunnelStatus): string {
		switch (status) {
			case 'running':
				return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
			case 'stopped':
				return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200';
			case 'starting':
			case 'stopping':
				return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200';
			case 'error':
				return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200';
			default:
				return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200';
		}
	}

	function getHealthColor(health: HealthStatus): string {
		switch (health) {
			case 'healthy':
				return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
			case 'unhealthy':
				return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200';
			case 'unknown':
			default:
				return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200';
		}
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold">Cloudflare Tunnels</h1>
			<p class="text-muted-foreground">
				Manage cloudflared tunnels for secure access to your services.
			</p>
		</div>
		<Button onclick={() => (showCreateDialog = true)}>
			<Plus class="h-4 w-4 mr-2" />
			Create Tunnel
		</Button>
	</div>

	<Card.Root>
		<Card.Header>
			<Card.Title>Active Tunnels</Card.Title>
			<Card.Description>View and manage your Cloudflare tunnels.</Card.Description>
		</Card.Header>
		<Card.Content>
			{#if tunnelsQuery.isLoading}
				<p class="text-sm text-muted-foreground text-center py-8">Loading tunnels...</p>
			{:else if tunnelsQuery.isError}
				<p class="text-sm text-destructive text-center py-8">
					Failed to load tunnels: {tunnelsQuery.error?.message}
				</p>
			{:else if !tunnelsQuery.data || tunnelsQuery.data.length === 0}
				<div class="text-center py-12">
					<Network class="h-12 w-12 mx-auto text-muted-foreground mb-4" />
					<p class="text-sm text-muted-foreground mb-4">No tunnels configured yet.</p>
					<Button onclick={() => (showCreateDialog = true)}>
						<Plus class="h-4 w-4 mr-2" />
						Create Your First Tunnel
					</Button>
				</div>
			{:else}
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>Name</Table.Head>
							<Table.Head>Status</Table.Head>
							<Table.Head>Health</Table.Head>
							<Table.Head class="text-right">Actions</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each tunnelsQuery.data as tunnel}
							<Table.Row>
								<Table.Cell class="font-medium">{tunnel.name}</Table.Cell>
								<Table.Cell>
									<Badge class={getStatusColor(tunnel.status)}>
										{tunnel.status}
									</Badge>
								</Table.Cell>
								<Table.Cell>
									<Badge class={getHealthColor(tunnel.health_status)}>
										{tunnel.health_status}
									</Badge>
								</Table.Cell>
								<Table.Cell class="text-right">
									<div class="flex justify-end gap-2">
										{#if tunnel.status === 'stopped' || tunnel.status === 'error'}
											<Button
												variant="ghost"
												size="sm"
												onclick={() => handleStart(tunnel)}
												title="Start tunnel"
											>
												<Play class="h-4 w-4" />
											</Button>
										{/if}
										{#if tunnel.status === 'running'}
											<Button
												variant="ghost"
												size="sm"
												onclick={() => handleStop(tunnel)}
												title="Stop tunnel"
											>
												<Square class="h-4 w-4" />
											</Button>
										{/if}
										{#if tunnel.status === 'running' || tunnel.status === 'error'}
											<Button
												variant="ghost"
												size="sm"
												onclick={() => handleRestart(tunnel)}
												title="Restart tunnel"
											>
												<RotateCw class="h-4 w-4" />
											</Button>
										{/if}
										<Button
											variant="ghost"
											size="sm"
											onclick={() => handleDelete(tunnel)}
											title="Delete tunnel"
										>
											<Trash2 class="h-4 w-4 text-destructive" />
										</Button>
									</div>
								</Table.Cell>
							</Table.Row>
							{#if tunnel.last_error}
								<Table.Row>
									<Table.Cell colspan={4} class="py-2">
										<div class="text-xs text-destructive bg-destructive/10 rounded p-2">
											Error: {tunnel.last_error}
										</div>
									</Table.Cell>
								</Table.Row>
							{/if}
						{/each}
					</Table.Body>
				</Table.Root>
			{/if}
		</Card.Content>
	</Card.Root>
</div>

<Dialog.Root bind:open={showCreateDialog}>
	<Dialog.Content class="sm:max-w-[500px]">
		<Dialog.Header>
			<Dialog.Title>Create Cloudflare Tunnel</Dialog.Title>
			<Dialog.Description>
				Create a cloudflared tunnel. Configure routing in the Cloudflare Zero Trust dashboard.
			</Dialog.Description>
		</Dialog.Header>
		<div class="space-y-4 py-4">
			<div class="space-y-2">
				<Label for="name">Name</Label>
				<Input id="name" bind:value={name} placeholder="my-tunnel" />
				<p class="text-xs text-muted-foreground">A friendly name for your tunnel.</p>
			</div>

			<div class="space-y-2">
				<Label for="token">Tunnel Token</Label>
				<Input id="token" type="password" bind:value={tunnelToken} placeholder="eyJh..." />
				<p class="text-xs text-muted-foreground">
					Your Cloudflare tunnel token from the Zero Trust dashboard.
				</p>
			</div>

			<div class="space-y-2">
				<Label for="project-id">Project ID (Optional)</Label>
				<Input id="project-id" bind:value={projectId} placeholder="project-uuid" />
				<p class="text-xs text-muted-foreground">
					Associate this tunnel with a specific project (optional).
				</p>
			</div>
		</div>
		<Dialog.Footer>
			<Button variant="outline" onclick={() => (showCreateDialog = false)}>Cancel</Button>
			<Button onclick={handleCreate} disabled={createTunnelMutation.isPending}>
				{createTunnelMutation.isPending ? 'Creating...' : 'Create Tunnel'}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
