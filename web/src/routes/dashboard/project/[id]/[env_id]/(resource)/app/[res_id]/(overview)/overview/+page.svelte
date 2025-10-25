<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Frame, FramePanel, FrameHeader, FrameTitle } from '$lib/components/ui/frame';
	import { Separator } from '$lib/components/ui/separator';
	import { toast } from 'svelte-sonner';
	import {
		Play,
		Square,
		RefreshCw,
		ExternalLink,
		GitBranch,
		GitCommit,
		FileText,
		RotateCcw,
		Github,
		Globe,
		ChevronDown
	} from 'lucide-svelte';
	import { createApplicationFetchQuery } from '$lib/features/applications/queries';
	import { createDeploymentsListQuery } from '$lib/features/deployments/queries';
	import {
		restartApplicationMutationQuery,
		startApplicationMutationQuery,
		stopApplicationMutationQuery
	} from '$lib/features/applications/mutations';

	const projectId = page.params.id!;
	const resId = page.params.res_id!;
	const envId = page.params.env_id!;

	const applicationQuery = createApplicationFetchQuery(projectId, envId, resId);
	const deploymentsQuery = createDeploymentsListQuery(projectId, resId);

	const application = $derived(applicationQuery.data);
	const latestDeployment = $derived(deploymentsQuery.data?.[0]);

	const startMutation = startApplicationMutationQuery(
		{
			projectId: projectId,
			environmentId: envId,
			resourceID: resId
		},
		{
			onSuccess: () => {
				toast.success('Application started successfully');
			},
			onError: (error: Error) => {
				toast.error(`Failed to start application: ${error.message}`);
			}
		}
	);

	const stopMutation = stopApplicationMutationQuery(
		{
			projectId: projectId,
			environmentId: envId,
			resourceID: resId
		},
		{
			onSuccess: () => {
				toast.success('Application stopped successfully');
			},
			onError: (error: Error) => {
				toast.error(`Failed to stop application: ${error.message}`);
			}
		}
	);

	const restartMutation = restartApplicationMutationQuery(
		{
			projectId: projectId,
			environmentId: envId,
			resourceID: resId
		},
		{
			onSuccess: () => {
				toast.success('Application restarted successfully');
			},
			onError: (error: Error) => {
				toast.error(`Failed to restart application: ${error.message}`);
			}
		}
	);

	const isAnyActionPending = $derived(
		startMutation.isPending || stopMutation.isPending || restartMutation.isPending
	);

	const isGitDeployment = $derived(
		application?.deployment_source?.source_type === 'git' ||
			application?.deployment_source?.type === 'git'
	);

	const gitUrl = $derived(
		application?.deployment_source?.git_url || application?.deployment_source?.git_repo?.url || null
	);

	const hasConfiguredDomain = $derived(
		!!(application?.custom_domain || application?.generated_domain || application?.domain)
	);

	const primaryDomain = $derived(
		application?.custom_domain || application?.generated_domain || application?.domain || null
	);

	function getStatusBadgeVariant(
		status: string
	): 'default' | 'secondary' | 'destructive' | 'outline' {
		switch (status) {
			case 'running':
			case 'success':
				return 'default';
			case 'stopped':
				return 'secondary';
			case 'failed':
				return 'destructive';
			case 'building':
			case 'deploying':
				return 'outline';
			default:
				return 'outline';
		}
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function formatRelativeTime(timestamp: string): string {
		const date = new Date(timestamp);
		const now = new Date();
		const diff = now.getTime() - date.getTime();
		const seconds = Math.floor(diff / 1000);
		const minutes = Math.floor(seconds / 60);
		const hours = Math.floor(minutes / 60);
		const days = Math.floor(hours / 24);

		if (days > 0) return `${days} day${days > 1 ? 's' : ''} ago`;
		if (hours > 0) return `${hours} hour${hours > 1 ? 's' : ''} ago`;
		if (minutes > 0) return `${minutes} minute${minutes > 1 ? 's' : ''} ago`;
		return `${seconds} second${seconds > 1 ? 's' : ''} ago`;
	}
</script>

{#if application}
	<div class="container mx-auto space-y-6 p-6">
		<div class="flex items-center justify-between">
			<h1 class="text-4xl font-bold">{application.name}</h1>
			<div class="flex items-center gap-2">
				{#if isGitDeployment && gitUrl}
					<Button size="lg" variant="outline" href={gitUrl} target="_blank">
						<Github class="mr-2 h-4 w-4" />
						Repository
					</Button>
				{/if}
				{#if application.status === 'stopped' || application.status === 'pending' || application.status === 'created'}
					<Button size="lg" disabled={isAnyActionPending} onclick={() => startMutation.mutate()}>
						<Play class="mr-2 h-4 w-4" />
						{startMutation.isPending ? 'Starting...' : 'Start'}
					</Button>
				{:else if application.status === 'running'}
					<Button
						size="lg"
						variant="outline"
						disabled={isAnyActionPending}
						onclick={() => stopMutation.mutate()}
					>
						<Square class="mr-2 h-4 w-4" />
						{stopMutation.isPending ? 'Stopping...' : 'Stop'}
					</Button>
				{/if}
				{#if application.status !== 'pending'}
					<Button
						size="lg"
						variant="outline"
						disabled={isAnyActionPending}
						onclick={() => restartMutation.mutate()}
					>
						<RefreshCw class="mr-2 h-4 w-4 {restartMutation.isPending ? 'animate-spin' : ''}" />
						{restartMutation.isPending ? 'Restarting...' : 'Restart'}
					</Button>
				{/if}
				<Button
					size="lg"
					variant="outline"
					href="/dashboard/project/{projectId}/{envId}/app/{resId}/domains"
				>
					<Globe class="mr-2 h-4 w-4" />
					Domains
				</Button>
				{#if hasConfiguredDomain && primaryDomain}
					<Button variant="outline" href="https://{primaryDomain}" target="_blank">
						Visit
						<ChevronDown class="ml-2 h-4 w-4" />
					</Button>
				{/if}
			</div>
		</div>

		<Separator class="-mx-48 data-[orientation=horizontal]:w-[calc(100%+21rem)]" />
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-3">
				<h1 class="text-xl font-semibold">Production Deployment</h1>
			</div>
			<div class="flex items-center gap-2">
				{#if latestDeployment}
					<Button
						size="sm"
						variant="outline"
						href="/dashboard/project/{projectId}/{envId}/app/{resId}/deployments/{latestDeployment.id}"
					>
						Build Logs
					</Button>
				{/if}
				<Button
					size="sm"
					variant="outline"
					href="/dashboard/project/{projectId}/{envId}/app/{resId}/logs"
				>
					Runtime Logs
				</Button>
				<Button size="sm" variant="outline" disabled>
					<RotateCcw class="mr-2 h-4 w-4" />
					Instant Rollback
				</Button>
			</div>
		</div>
		<Frame>
			<FramePanel>
				{#if latestDeployment}
					<div class="px-5 py-4 space-y-4">
						<div class="grid grid-cols-2 gap-6">
							<div>
								<p class="text-sm text-muted-foreground mb-1">Deployment</p>
								<p class="font-medium text-sm break-all">
									{latestDeployment.id || 'N/A'}
								</p>
							</div>
							<div>
								<p class="text-sm text-muted-foreground mb-1">Domains</p>
								{#if hasConfiguredDomain && primaryDomain}
									<a
										href="https://{primaryDomain}"
										target="_blank"
										class="font-medium text-sm hover:underline flex items-center gap-1"
									>
										{primaryDomain}
										<ExternalLink class="h-3 w-3" />
									</a>
								{:else}
									<p class="text-sm text-muted-foreground">No domain configured</p>
								{/if}
							</div>
						</div>

						<div class="grid grid-cols-2 gap-6">
							<div>
								<p class="text-sm text-muted-foreground mb-1">Status</p>
								<div class="flex items-center gap-2">
									<Badge variant={getStatusBadgeVariant(latestDeployment.status)}>
										{latestDeployment.status}
									</Badge>
									<span class="text-sm">
										{formatRelativeTime(latestDeployment.created_at)}
									</span>
								</div>
							</div>
							<div>
								<p class="text-sm text-muted-foreground mb-1">Created</p>
								<p class="text-sm">
									{formatDate(latestDeployment.created_at)}
									{#if latestDeployment.author}
										by {latestDeployment.author}
									{/if}
								</p>
							</div>
						</div>

						{#if latestDeployment.branch || latestDeployment.commit_hash}
							<div class="grid grid-cols-2 gap-6">
								{#if latestDeployment.branch}
									<div>
										<div class="flex items-center gap-2 text-sm text-muted-foreground mb-1">
											<GitBranch class="h-4 w-4" />
											<span>Source</span>
										</div>
										<p class="font-medium flex items-center gap-1">
											<GitBranch class="h-4 w-4" />
											{latestDeployment.branch}
										</p>
									</div>
								{/if}
								{#if latestDeployment.commit_hash}
									<div>
										<div class="flex items-center gap-2 text-sm text-muted-foreground mb-1">
											<GitCommit class="h-4 w-4" />
											<span>Commit</span>
										</div>
										<p class="font-mono font-medium flex items-center gap-1">
											<GitCommit class="h-4 w-4" />
											{latestDeployment.commit_hash.slice(0, 7)}
											{latestDeployment.commit_message || ''}
										</p>
									</div>
								{/if}
							</div>
						{/if}

						{#if application.deployment_source}
							{@const source = application.deployment_source}
							{#if source.source_type === 'git' || source.type === 'git'}
								<div class="grid grid-cols-2 gap-6">
									{#if source.git_branch || source.git_repo?.branch}
										<div>
											<span class=" text-muted-foreground mb-1 text-sm"> Source </span>
											<a
												class="font-medium flex items-center gap-2 text-sm"
												href={source.git_url || source.git_repo?.url}
											>
												<GitBranch class="h-4 w-4" />
												<span>{source.git_branch || source.git_repo?.branch}</span>
											</a>
										</div>
									{/if}
									{#if source.git_path || source.git_repo?.path}
										<div>
											<div class="flex items-center gap-2 text-sm text-muted-foreground mb-1">
												<span>Path</span>
											</div>
											<p
												class="font-medium font-mono flex items-center gap-2 text-sm text-muted-foreground"
											>
												<span>{source.git_path || source.git_repo?.path}</span>
											</p>
										</div>
									{/if}
								</div>
							{:else if source.source_type === 'docker' || source.type === 'registry'}
								<div>
									<div class="flex items-center gap-2 text-sm text-muted-foreground mb-1">
										<span>Docker Image</span>
									</div>
									<p class="font-medium font-mono text-sm">
										{source.docker_image || source.registry?.image || 'Not configured'}
									</p>
								</div>
							{/if}
						{:else}
							<p class="text-muted-foreground">No source configured</p>
						{/if}
					</div>
				{:else}
					<div class="text-center py-8 text-muted-foreground px-5">
						<p>No deployments yet</p>
					</div>
				{/if}
			</FramePanel>
		</Frame>
	</div>
{/if}
