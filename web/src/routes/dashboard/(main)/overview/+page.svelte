<script lang="ts">
	import { goto } from '$app/navigation';
	import { formatTimeAgo } from '$lib/utils/dates';

	import { Button } from '$lib/components/ui/button';
	import { Ellipsis, Globe, Plus } from 'lucide-svelte';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import * as Empty from '$lib/components/ui/empty/index.js';

	import { createProjectsQuery } from '$lib/features/projects/queries';
	import { createActivitiesQuery } from '$lib/features/activities/queries';
	import { createServersListQuery } from '$lib/features/servers/queries';

	const projects = createProjectsQuery();
	const activitiesQuery = createActivitiesQuery(10);
	const serversQuery = createServersListQuery();

	function goToProject(projectId: string) {
		goto(`/dashboard/project/${projectId}`);
	}

	function getActivityIcon(action: string) {
		if (action.includes('create')) return '✓';
		if (action.includes('delete')) return '✕';
		if (action.includes('update')) return '↻';
		if (action.includes('deploy')) return '🚀';
		return '•';
	}

	function getActivityColor(action: string) {
		if (action.includes('create')) return 'bg-success';
		if (action.includes('delete')) return 'bg-destructive';
		if (action.includes('update')) return 'bg-info';
		if (action.includes('deploy')) return 'bg-accent';
		return 'bg-muted';
	}

	function getServerStatusColor(status: string) {
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
	<title>Dashboard - mikrocloud</title>
</svelte:head>

<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
	<div class="lg:col-span-2 space-y-6">
		<div>
			<h2 class="text-lg font-semibold mb-4">Projects</h2>
			{#if projects.isLoading}
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					{#each Array(4) as _}
						<Skeleton class="h-32 w-full rounded-lg" />
					{/each}
				</div>
			{:else if projects.data && projects.data.length === 0}
				<Empty.Root class="border border-border">
					<Empty.Header>
						<Empty.Media variant="icon">
							<Plus />
						</Empty.Media>
						<Empty.Title>No Projects Yet</Empty.Title>
						<Empty.Description>
							You haven't created any projects yet. Get started by creating your first project.
						</Empty.Description>
					</Empty.Header>
					<Empty.Content>
						<div>
							<Button href="/dashboard/projects/new">Create Project</Button>
						</div>
					</Empty.Content>
				</Empty.Root>
			{:else if projects.data}
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					{#each projects.data as project (project.id)}
						<div
							class="bg-card border border-border rounded-lg p-4 hover:bg-card-hover transition-colors cursor-pointer"
							onclick={() => goToProject(project.id)}
							role="button"
							tabindex="0"
							onkeydown={(e) => e.key === 'Enter' && goToProject(project.id)}
						>
							<div class="flex items-center justify-between mb-8">
								<div class="flex items-center space-x-2">
									<Globe class="w-4 h-4 text-gray-400" />
									<span class="text-sm">{project.name}</span>
								</div>
								<div
									class="text-gray-400 hover:text-white cursor-pointer"
									onclick={(e) => e.stopPropagation()}
									role="button"
									tabindex="0"
									onkeydown={(e) => e.key === 'Enter' && e.stopPropagation()}
								>
									<Ellipsis class="w-5 h-5" />
								</div>
							</div>
							<div class="flex items-center justify-between">
								<span class="text-xs text-gray-500">
									Created {formatTimeAgo(project.created_at)}
								</span>
								{#if project.description}
									<span class="text-xs text-gray-400 truncate max-w-[150px]">
										{project.description}
									</span>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>

		<div>
			<h2 class="text-lg font-semibold mb-4">Servers</h2>
			<div class="bg-card border border-border rounded-lg p-4">
				{#if serversQuery.isLoading}
					<div class="space-y-4">
						<Skeleton class="h-24 w-full" />
						<Skeleton class="h-24 w-full" />
					</div>
				{:else if serversQuery.data && serversQuery.data.length === 0}
					<div class="text-center text-gray-400 text-sm py-4">No servers configured</div>
				{:else if serversQuery.data}
					{#each serversQuery.data as server (server.id)}
						<div class="pb-4 mb-4 last:pb-0 last:mb-0 border-b border-border last:border-0">
							<div class="flex items-center justify-between mb-2">
								<div class="flex items-center space-x-2">
									<span class="text-sm">🖥️</span>
									<span class="text-sm font-medium">{server.name}</span>
								</div>
								<div class="flex items-center space-x-2">
									<div class="w-2 h-2 rounded-full {getServerStatusColor(server.status)}"></div>
									<span class="text-xs text-gray-400">{server.status}</span>
								</div>
							</div>
							<div class="text-xs text-gray-500 mb-2">{server.hostname}</div>
							{#if server.tags && server.tags.length > 0}
								<div class="flex flex-wrap gap-1">
									{#each server.tags as tag}
										<span class="text-xs bg-card border border-border rounded px-2 py-0.5">
											{tag}
										</span>
									{/each}
								</div>
							{/if}
						</div>
					{/each}
				{/if}
			</div>
		</div>
	</div>

	<div>
		<h2 class="text-lg font-semibold mb-4">Activity</h2>
		<div class="bg-white/5 border border-white/10 rounded-lg p-4">
			{#if activitiesQuery.isLoading}
				<div class="space-y-6">
					{#each Array(5) as _}
						<div class="flex items-start space-x-3">
							<Skeleton class="h-8 w-8 rounded-full" />
							<div class="flex-1 space-y-2">
								<Skeleton class="h-4 w-32" />
								<Skeleton class="h-3 w-48" />
								<Skeleton class="h-5 w-24" />
							</div>
						</div>
					{/each}
				</div>
			{:else if activitiesQuery.data && activitiesQuery.data.activities.length === 0}
				<Empty.Root>
					<Empty.Header>
						<Empty.Title>No recent activity</Empty.Title>
					</Empty.Header>
				</Empty.Root>
			{:else if activitiesQuery.data}
				<div class="space-y-6">
					{#each activitiesQuery.data.activities as activity, index (activity.id)}
						<div class="flex items-start space-x-3">
							<div class="relative">
								<div
									class="w-8 h-8 rounded-full {getActivityColor(
										activity.activity_type
									)} flex items-center justify-center text-xs shrink-0"
								>
									{getActivityIcon(activity.activity_type)}
								</div>
								{#if index !== activitiesQuery.data.activities.length - 1}
									<div class="absolute top-8 left-4 w-px h-8 bg-border"></div>
								{/if}
							</div>
							<div class="flex-1 min-w-0">
								<div class="flex items-center justify-between mb-1">
									<span class="text-sm font-medium">{activity.initiator_name || 'System'}</span>
									<span class="text-xs text-gray-500">{formatTimeAgo(activity.created_at)}</span>
								</div>
								<div class="text-xs text-gray-400 mb-1">{activity.activity_type}</div>
								{#if activity.resource_type && activity.resource_id}
									<div class="text-xs bg-card border border-border rounded px-2 py-1 inline-block">
										{activity.resource_type}: {activity.resource_id.substring(0, 8)}
									</div>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>
