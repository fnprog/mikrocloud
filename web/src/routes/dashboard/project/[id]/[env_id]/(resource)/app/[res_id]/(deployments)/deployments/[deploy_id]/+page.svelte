<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { deploymentsApi, type DeploymentStatus } from '$lib/features/deployments/api';
	import type { StructuredLog } from '$lib/features/deployments/types';
	import LogViewer from '$lib/features/deployments/components/LogViewer.svelte';
	import { toast } from 'svelte-sonner';
	import { onDestroy } from 'svelte';
	import {
		CheckCircle,
		AlertCircle,
		XCircle,
		Clock,
		GitBranch,
		GitCommit,
		User,
		RotateCcw,
		ArrowLeft,
		PlayCircle,
		StopCircle
	} from 'lucide-svelte';
	import { formatTimeAgo } from '$lib/utils/dates';

	const projectId = $derived(page.params.id);
	const envId = $derived(page.params.env_id);
	const resId = $derived(page.params.res_id);
	const deployId = $derived(page.params.deploy_id);

	const queryClient = useQueryClient();

	let streamedLogs = $state<StructuredLog[]>([]);
	let isStreaming = $state(false);
	let closeStream: (() => void) | null = null;

	const deploymentQuery = createQuery(() => ({
		queryKey: ['deployment', projectId, resId, deployId],
		queryFn: () => deploymentsApi.get(projectId, resId, deployId),
		enabled: !!projectId && !!resId && !!deployId,
		refetchInterval: (query) => {
			const data = query.state.data;
			if (data && ['pending', 'building', 'deploying'].includes(data.status) && !isStreaming) {
				return 3000;
			}
			return false;
		}
	}));

	const deployment = $derived(deploymentQuery.data);

	const logsQuery = createQuery(() => ({
		queryKey: ['deployment-logs', projectId, resId, deployId],
		queryFn: () => deploymentsApi.getLogs(projectId, resId, deployId, 'build'),
		enabled:
			!!deployment &&
			!['pending', 'building', 'deploying'].includes(deployment.status) &&
			!isStreaming
	}));

	const completedLogs = $derived(logsQuery.data);

	$effect(() => {
		if (
			deployment &&
			['pending', 'building', 'deploying'].includes(deployment.status) &&
			!isStreaming
		) {
			startLogStream();
		} else if (
			deployment &&
			!['pending', 'building', 'deploying'].includes(deployment.status) &&
			isStreaming
		) {
			stopLogStream();
		}
	});

	function startLogStream() {
		if (isStreaming) return;

		isStreaming = true;
		streamedLogs = [];

		closeStream = deploymentsApi.streamLogs(
			projectId,
			resId,
			deployId,
			(log: StructuredLog) => {
				streamedLogs = [...streamedLogs, log];
			},
			(status: DeploymentStatus) => {
				isStreaming = false;
				queryClient.invalidateQueries({ queryKey: ['deployment', projectId, resId, deployId] });
			},
			(error: Error) => {
				console.error('Log stream error:', error);
				isStreaming = false;
			}
		);
	}

	function stopLogStream() {
		if (closeStream) {
			closeStream();
			closeStream = null;
		}
		isStreaming = false;
	}

	onDestroy(() => {
		stopLogStream();
	});

	const redeployMutation = createMutation(() => ({
		mutationFn: () => deploymentsApi.redeploy(projectId, resId, deployId),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['deployment', projectId, resId, deployId] });
			queryClient.invalidateQueries({ queryKey: ['deployments', projectId, resId] });
			toast.success('Redeployment started');
		},
		onError: (error: Error) => {
			toast.error(`Failed to redeploy: ${error.message}`);
		}
	}));

	const cancelMutation = createMutation(() => ({
		mutationFn: () => deploymentsApi.cancel(projectId, resId, deployId),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['deployment', projectId, resId, deployId] });
			queryClient.invalidateQueries({ queryKey: ['deployments', projectId, resId] });
			toast.success('Deployment cancelled');
		},
		onError: (error: Error) => {
			toast.error(`Failed to cancel deployment: ${error.message}`);
		}
	}));

	function getStatusIcon(status: DeploymentStatus) {
		switch (status) {
			case 'success':
				return CheckCircle;
			case 'failed':
				return AlertCircle;
			case 'cancelled':
				return XCircle;
			case 'building':
			case 'deploying':
			case 'pending':
				return Clock;
			default:
				return Clock;
		}
	}

	function getStatusColor(status: DeploymentStatus) {
		switch (status) {
			case 'success':
				return 'text-green-500';
			case 'failed':
				return 'text-red-500';
			case 'cancelled':
				return 'text-gray-500';
			case 'building':
			case 'deploying':
			case 'pending':
				return 'text-blue-500';
			default:
				return 'text-gray-500';
		}
	}

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

	function formatDuration(seconds?: number): string {
		if (!seconds) return 'N/A';
		const mins = Math.floor(seconds / 60);
		const secs = seconds % 60;
		return `${mins}m ${secs}s`;
	}

	function formatDateTime(timestamp: string): string {
		const date = new Date(timestamp);
		return date.toLocaleString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit'
		});
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

	function goBack() {
		goto(`/dashboard/project/${projectId}/${envId}/app/${resId}/deployments`);
	}

	function handleRedeploy() {
		redeployMutation.mutate();
	}

	function handleCancel() {
		cancelMutation.mutate();
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div class="flex items-center space-x-4">
			<Button variant="ghost" size="sm" onclick={goBack}>
				<ArrowLeft class="w-4 h-4 mr-2" />
				Back
			</Button>
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Deployment Details</h2>
				<p class="text-muted-foreground">View detailed information about this deployment</p>
			</div>
		</div>

		{#if deployment}
			<div class="flex items-center space-x-2">
				{#if ['pending', 'building', 'deploying'].includes(deployment.status)}
					<Button
						size="sm"
						variant="outline"
						disabled={cancelMutation.isPending}
						onclick={handleCancel}
					>
						<StopCircle class="w-4 h-4 mr-2" />
						Cancel
					</Button>
				{/if}
				{#if deployment.status === 'success'}
					<Button
						size="sm"
						variant="outline"
						disabled={redeployMutation.isPending}
						onclick={handleRedeploy}
					>
						<RotateCcw class="w-4 h-4 mr-2 {redeployMutation.isPending ? 'animate-spin' : ''}" />
						Redeploy
					</Button>
				{/if}
			</div>
		{/if}
	</div>

	{#if deploymentQuery.isLoading}
		<div class="flex items-center justify-center py-12">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
		</div>
	{:else if deploymentQuery.isError}
		<Card>
			<CardContent class="p-6">
				<div class="text-center">
					<AlertCircle class="mx-auto h-12 w-12 text-destructive mb-4" />
					<h3 class="text-lg font-medium">Failed to load deployment</h3>
					<p class="text-muted-foreground mt-2">
						{deploymentQuery.error?.message || 'An error occurred'}
					</p>
					<Button variant="outline" class="mt-4" onclick={goBack}>Go Back</Button>
				</div>
			</CardContent>
		</Card>
	{:else if deployment}
		<div class="grid gap-6">
			<div class="bg-card/80 rounded-xl">
				<div class="flex px-6 py-6 justify-between text-muted-foreground">
					<div>
						<p>Created</p>
						<div class="flex gap-2">
							{deployment.triggered_by_username}
							{formatTimeAgo(deployment.created_at)}
						</div>
					</div>

					<div>
						<p>Status</p>
						<div class="flex items-center gap-2">
							<span class="w-2 h-2 rounded-full {getStatusDotColor(deployment.status)}"></span>
							<span class="text-sm font-medium">
								{getStatusText(deployment.status)}
							</span>
						</div>
					</div>

					<div>
						<p>Time to Ready</p>
						<div class="flex">
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
					</div>
				</div>

				<Card class="mx-1 rounded-b-none">
					<CardHeader>
						<CardTitle>Domains</CardTitle>
					</CardHeader>
					<CardContent>
						<div class="flex flex-col items-start space-y-3">
							<p>example.com</p>
							<p>asfsafass.example.com</p>
							<p>sasfsfdsasa.example.com</p>
						</div>
					</CardContent>
				</Card>

				<Card class="mx-1 rounded-t-none">
					<CardHeader>
						<CardTitle>Source</CardTitle>
					</CardHeader>
					<CardContent>
						<div class="space-y-4">
							{#if deployment.branch}
								<div>
									<h4 class="text-sm font-medium text-muted-foreground mb-1 flex items-center">
										<GitBranch class="w-4 h-4 mr-1" />
										Branch
									</h4>
									<p class="text-base">{deployment.branch}</p>
								</div>
							{/if}

							<div class="grid grid-cols-2 gap-4">
								{#if deployment.commit_hash}
									<div>
										<h4 class="text-sm font-medium text-muted-foreground mb-1 flex items-center">
											<GitCommit class="w-4 h-4 mr-1" />
											Commit Hash
										</h4>
										<p class="text-base font-mono">{deployment.commit_hash}</p>
									</div>
								{/if}

								{#if deployment.commit_message}
									<div>
										<h4 class="text-sm font-medium text-muted-foreground mb-1">Message</h4>
										<p class="text-base">{deployment.commit_message}</p>
									</div>
								{/if}

								<!-- {#if deployment.author} -->
								<!-- 	<div> -->
								<!-- 		<h4 class="text-sm font-medium text-muted-foreground mb-1 flex items-center"> -->
								<!-- 			<User class="w-4 h-4 mr-1" /> -->
								<!-- 			Author -->
								<!-- 		</h4> -->
								<!-- 		<p class="text-base">{deployment.author}</p> -->
								<!-- 	</div> -->
								<!-- {/if} -->
							</div>
						</div>
					</CardContent>
				</Card>
			</div>

			<Card>
				<CardHeader>
					<CardTitle class="flex items-center justify-between">
						<span>Build Logs</span>
						{#if isStreaming}
							<Badge variant="outline" class="text-xs">
								<span class="inline-block w-2 h-2 bg-green-500 rounded-full mr-2 animate-pulse"
								></span>
								Live
							</Badge>
						{/if}
					</CardTitle>
				</CardHeader>
				<CardContent>
					{#if streamedLogs.length > 0}
						<LogViewer logs={streamedLogs} {isStreaming} />
					{:else if Array.isArray(completedLogs)}
						<LogViewer logs={completedLogs} isStreaming={false} />
					{:else if typeof completedLogs === 'string' && completedLogs}
						<div class="bg-muted rounded-md p-4 overflow-x-auto overflow-y-auto max-h-[600px]">
							<pre class="text-sm font-mono whitespace-pre-wrap">{completedLogs}</pre>
						</div>
					{:else if deployment.build_logs}
						<div class="bg-muted rounded-md p-4 overflow-x-auto overflow-y-auto max-h-[600px]">
							<pre class="text-sm font-mono whitespace-pre-wrap">{deployment.build_logs}</pre>
						</div>
					{:else}
						<p class="text-sm text-muted-foreground">No logs available yet...</p>
					{/if}
				</CardContent>
			</Card>
		</div>
	{/if}
</div>
