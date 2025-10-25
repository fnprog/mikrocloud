<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import * as Frame from '$lib/components/ui/frame';
	import * as Collapsible from '$lib/components/ui/collapsible/index.js';
	import { Separator } from '$lib/components/ui/separator';
	import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { deploymentsApi, type DeploymentStatus } from '$lib/features/deployments/api';
	import type { StructuredLog } from '$lib/features/deployments/types';
	import LogViewer from '$lib/features/deployments/components/LogViewer.svelte';
	import { toast } from 'svelte-sonner';
	import { onDestroy } from 'svelte';
	import {
		AlertCircle,
		GitBranch,
		GitCommit,
		RotateCcw,
		ArrowLeft,
		StopCircle,
		ChevronDown
	} from 'lucide-svelte';
	import { formatTimeAgo } from '$lib/utils/dates';

	const projectId = $derived(page.params.id);
	const envId = $derived(page.params.env_id);
	const resId = $derived(page.params.res_id);
	const deployId = $derived(page.params.deploy_id);

	const queryClient = useQueryClient();

	let streamedLogs = $state<StructuredLog[]>([]);
	let isStreaming = $state(false);
	let isBuildLogOpen = $state(true);
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
		<div>
			<h2 class="text-2xl mb-3 font-bold tracking-tight">Deployment Details</h2>
			<p class="text-muted-foreground">View detailed information about this deployment</p>
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
			<Frame.Frame class="w-full">
				<Frame.FrameHeader>
					<div class="flex justify-between text-muted-foreground">
						<div class="space-y-1">
							<p class="font-semibold">Created</p>
							<div class="flex gap-2">
								{deployment.triggered_by_username}
								{formatTimeAgo(deployment.created_at)}
							</div>
						</div>

						<div class="space-y-1">
							<p class="font-semibold">Status</p>
							<div class="flex items-center gap-2 text-muted-foreground">
								<span class="w-2 h-2 rounded-full {getStatusDotColor(deployment.status)}"></span>
								<span class=" font-medium">
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
				</Frame.FrameHeader>
				<Frame.FramePanel class="space-y-4 ">
					<div>
						<h2 class="text-sm font-semibold text-muted-foreground mb-1">Domains</h2>
						<div class="flex flex-col items-start gap-2 text-sm">
							<p>example.com</p>
							<p>asfsafass.example.com</p>
							<p>sasfsfdsasa.example.com</p>
						</div>
					</div>
					<Separator class="-mx-5 data-[orientation=horizontal]:w-[calc(100%+2.5rem)] " />
					<div>
						<h2 class="text-sm font-semibold text-muted-foreground mb-1">Source</h2>
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
					</div>
				</Frame.FramePanel>
			</Frame.Frame>

			<Card class="pb-0">
				<Collapsible.Root bind:open={isBuildLogOpen}>
					<CardHeader class="pb-4">
						<CardTitle class="flex items-center justify-between">
							<div class="flex items-center gap-3">
								<Collapsible.Trigger>
									<ChevronDown />
									<span class="sr-only">Toggle</span>
								</Collapsible.Trigger>
								<span class="text-muted-foreground">Build Logs</span>
							</div>
							{#if isStreaming}
								<Badge variant="outline" class="text-xs">
									<span class="inline-block w-2 h-2 bg-green-500 rounded-full mr-2 animate-pulse"
									></span>
									Live
								</Badge>
							{/if}
						</CardTitle>
					</CardHeader>
					<Collapsible.Content class="p-1">
						<CardContent
							class="relative bg-clip-padding rounded-xl border bg-muted p-5 px-0 pt-0 shadow-xs before:pointer-events-none before:absolute before:inset-0 before:rounded-[calc(var(--radius-xl)-1px)] before:shadow-[0_1px_--theme(--color-black/4%)] has-[table]:before:hidden dark:bg-clip-border dark:before:shadow-[0_-1px_--theme(--color-white/8%)]"
						>
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
					</Collapsible.Content>
				</Collapsible.Root>
			</Card>
		</div>
	{/if}
</div>
