<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent } from '$lib/components/ui/card';
	import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { deploymentsApi, type DeploymentStatus } from '$lib/features/deployments/api';
	import { toast } from 'svelte-sonner';
	import { AlertCircle, GitBranch, GitCommit } from 'lucide-svelte';

	const projectId = $derived(page.params.id);
	const envId = $derived(page.params.env_id);
	const resId = $derived(page.params.res_id);

	const queryClient = useQueryClient();

	const deploymentsQuery = createQuery(() => ({
		queryKey: ['deployments', projectId, resId],
		queryFn: () => deploymentsApi.list(projectId, resId),
		enabled: !!projectId && !!resId,
		refetchInterval: 5000
	}));

	const deployments = $derived(deploymentsQuery.data || []);

	const redeployMutation = createMutation(() => ({
		mutationFn: (deploymentId: string) => deploymentsApi.redeploy(projectId, resId, deploymentId),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['deployments', projectId, resId] });
			toast.success('Redeployment started');
		},
		onError: (error: Error) => {
			toast.error(`Failed to redeploy: ${error.message}`);
		}
	}));

	function getStatusDotColor(status: DeploymentStatus): string {
		switch (status) {
			case 'running':
				return 'bg-green-500';
			case 'failed':
				return 'bg-red-500';
			case 'building':
			case 'deploying':
			case 'pending':
			case 'queued':
				return 'bg-orange-500';
			case 'cancelled':
			case 'stopped':
				return 'bg-gray-500';
			default:
				return 'bg-gray-500';
		}
	}

	function getStatusText(status: DeploymentStatus): string {
		switch (status) {
			case 'running':
				return 'Ready';
			case 'building':
				return 'Building';
			case 'deploying':
				return 'Deploying';
			case 'pending':
				return 'Pending';
			case 'queued':
				return 'Queued';
			case 'failed':
				return 'Failed';
			case 'cancelled':
				return 'Cancelled';
			case 'stopped':
				return 'Stopped';
			default:
				return status;
		}
	}

	function formatCompactDuration(seconds?: number): string {
		if (!seconds) return '';
		if (seconds < 60) return `${seconds}s`;
		const mins = Math.floor(seconds / 60);
		return `${mins}m`;
	}

	function formatCompactTime(timestamp: string): string {
		const date = new Date(timestamp);
		const now = new Date();
		const diff = now.getTime() - date.getTime();
		const seconds = Math.floor(diff / 1000);
		const minutes = Math.floor(seconds / 60);
		const hours = Math.floor(minutes / 60);
		const days = Math.floor(hours / 24);

		if (days > 0) return `${days}d ago`;
		if (hours > 0) return `${hours}h ago`;
		if (minutes > 0) return `${minutes}m ago`;
		return `${seconds}s ago`;
	}

	function viewDeployment(deploymentId: string) {
		goto(`/dashboard/project/${projectId}/${envId}/app/${resId}/deployments/${deploymentId}`);
	}

	function redeployCommit(deploymentId: string) {
		redeployMutation.mutate(deploymentId);
	}

	function isCurrentDeployment(status: DeploymentStatus): boolean {
		return status === 'running';
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h2 class="text-2xl font-bold tracking-tight">Deployments</h2>
			<p class="text-muted-foreground">View and manage all deployments for this application</p>
		</div>
	</div>

	{#if deploymentsQuery.isLoading}
		<div class="flex items-center justify-center py-12">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
		</div>
	{:else if deploymentsQuery.isError}
		<Card>
			<CardContent class="p-6">
				<div class="text-center">
					<AlertCircle class="mx-auto h-12 w-12 text-destructive mb-4" />
					<h3 class="text-lg font-medium">Failed to load deployments</h3>
					<p class="text-muted-foreground mt-2">
						{deploymentsQuery.error?.message || 'An error occurred'}
					</p>
				</div>
			</CardContent>
		</Card>
	{:else if deployments.length === 0}
		<Card>
			<CardContent class="p-12">
				<div class="text-center">
					<GitBranch class="mx-auto h-12 w-12 text-muted-foreground mb-4" />
					<h3 class="text-lg font-medium">No deployments yet</h3>
					<p class="text-muted-foreground mt-2">
						Deployments will appear here once you trigger your first deployment.
					</p>
				</div>
			</CardContent>
		</Card>
	{:else}
		<div class="space-y-2">
			{#each deployments as deployment (deployment.id)}
				<button
					class="w-full text-left border border-border rounded-lg p-4 hover:bg-accent/50 transition-colors cursor-pointer"
					onclick={() => viewDeployment(deployment.id)}
				>
					<div class="grid grid-cols-[1fr_auto_auto] gap-6 items-center">
						<div class="flex flex-col gap-1">
							<div class="flex items-center gap-2">
								<span class="font-mono text-sm font-medium">
									{deployment.id.slice(0, 8)}
								</span>
								{#if isCurrentDeployment(deployment.status)}
									<Badge variant="default" class="text-xs">Current</Badge>
								{/if}
							</div>
							<span class="text-xs text-muted-foreground">
								{deployment.is_production ? 'Production' : 'Preview'}
								{#if deployment.triggered_by_username}
									· by {deployment.triggered_by_username}
								{/if}
							</span>
						</div>

						<div class="flex flex-col items-end gap-1">
							<div class="flex items-center gap-2">
								<span class="w-2 h-2 rounded-full {getStatusDotColor(deployment.status)}"></span>
								<span class="text-sm font-medium">
									{getStatusText(deployment.status)}
								</span>
							</div>
							<div class="flex items-center gap-2 text-xs text-muted-foreground">
								{#if deployment.build_duration_seconds !== undefined || deployment.deploy_duration_seconds !== undefined}
									<span
										>{formatCompactDuration(
											(deployment.build_duration_seconds || 0) +
												(deployment.deploy_duration_seconds || 0)
										)}</span
									>
								{/if}
								<span>({formatCompactTime(deployment.created_at)})</span>
							</div>
						</div>

						<div class="flex items-center gap-4">
							{#if deployment.branch}
								<div class="flex items-center gap-1.5 text-sm">
									<GitBranch class="w-4 h-4 text-muted-foreground" />
									<span>{deployment.branch}</span>
								</div>
							{/if}
							{#if deployment.commit_hash}
								<div class="flex items-center gap-1.5 text-sm">
									<GitCommit class="w-4 h-4 text-muted-foreground" />
									<span class="font-mono">{deployment.commit_hash.slice(0, 7)}</span>
								</div>
							{/if}
						</div>
					</div>
				</button>
			{/each}
		</div>
	{/if}
</div>
