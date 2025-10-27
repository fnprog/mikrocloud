<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	import { Search } from 'lucide-svelte';
	import EnvironmentTabs from '$lib/components/projects/environment-tabs.svelte';
	import AddResourceMenu from '$lib/components/projects/add-resource-menu.svelte';
	import ApplicationCard from '$lib/components/projects/application-card.svelte';
	import DatabaseCard from '$lib/components/projects/database-card.svelte';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';

	import { createProjectQuery } from '$lib/features/projects/queries';
	import { createEnvironmentsListQuery } from '$lib/features/environments/queries';
	import { createApplicationsFetchQuery } from '$lib/features/applications/queries';
	import { createDatabasesFetchQuery } from '$lib/features/databases/queries';

	import type { Application } from '$lib/features/applications/types';
	import type { Database } from '$lib/features/databases/types';

	const projectId = page.params.id!;

	let selectedEnvironmentId = $state<string | undefined>(undefined);
	let searchQuery = $state('');

	const projectQuery = createProjectQuery(projectId);
	const environmentsQuery = createEnvironmentsListQuery(projectId);
	const applicationsQuery = createApplicationsFetchQuery(projectId, '');
	const databasesQuery = createDatabasesFetchQuery(projectId, '');

	$effect(() => {
		if (environmentsQuery.data && !selectedEnvironmentId) {
			const productionEnv = environmentsQuery.data.find((env) => env.name === 'production');
			if (productionEnv) {
				selectedEnvironmentId = productionEnv.id;
			}
		}
	});

	const filteredResources = $derived.by(() => {
		const apps: Application[] = applicationsQuery.data || [];
		const dbs: Database[] = databasesQuery.data || [];

		let filtered: Array<{ type: 'application' | 'database'; data: Application | Database }> = [
			...apps.map((app) => ({ type: 'application' as const, data: app })),
			...dbs.map((db) => ({ type: 'database' as const, data: db }))
		];

		if (selectedEnvironmentId) {
			filtered = filtered.filter(
				(item) =>
					('environment_id' in item.data && item.data.environment_id === selectedEnvironmentId) ||
					('environment' in item.data && item.data.environment === selectedEnvironmentId)
			);
		}

		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase();
			filtered = filtered.filter((item) => item.data.name.toLowerCase().includes(query));
		}

		return filtered;
	});

	const resourceCounts = $derived.by(() => {
		const apps = applicationsQuery.data || [];
		const dbs = databasesQuery.data || [];
		const envs = environmentsQuery.data || [];

		const counts: Record<string, number> = {
			all: apps.length + dbs.length
		};

		envs.forEach((env) => {
			const envApps = apps.filter((app) => app.environment_id === env.id);
			const envDbs = dbs.filter((db) => 'environment' in db && db.environment === env.id);
			counts[env.id] = envApps.length + envDbs.length;
		});

		return counts;
	});

	function handleAddEnvironment() {
		console.log('Add environment clicked');
	}
</script>

<!-- TODO: Better loading and skeletton for the active stuff -->
<div class="flex flex-col gap-6 max-w-7xl mx-auto p-6">
	{#if projectQuery.isLoading}
		<div class="text-muted-foreground">Loading project...</div>
	{:else if projectQuery.error}
		<div class="text-destructive">Error loading project: {projectQuery.error.message}</div>
	{:else if projectQuery.data}
		<div class="space-y-6">
			<div class="my-[40px] w-full mx-auto">
				<h1 class="text-3xl font-bold">{projectQuery.data.name}</h1>
				{#if projectQuery.data.description}
					<p class="text-muted-foreground mt-4">{projectQuery.data.description}</p>
				{/if}
			</div>
			<div class="mt-[46px]"></div>

			<!-- TODO: Add a modal to create new environments -->
			<EnvironmentTabs
				environments={environmentsQuery.data || []}
				bind:selectedEnvironmentId
				onSelect={(envId) => (selectedEnvironmentId = envId)}
				onAdd={handleAddEnvironment}
				counts={resourceCounts}
			/>

			<div class="flex items-center justify-between gap-4">
				<InputGroup.Root class=" max-w-md">
					<InputGroup.Input placeholder="Search projects..." bind:value={searchQuery} />
					<InputGroup.Addon>
						<Search />
					</InputGroup.Addon>
				</InputGroup.Root>

				<AddResourceMenu {projectId} envId={selectedEnvironmentId} />
			</div>

			{#if applicationsQuery.isLoading || databasesQuery.isLoading}
				<div class="text-muted-foreground">Loading resources...</div>
			{:else if filteredResources.length === 0}
				<Empty.Root>
					<Empty.Header>
						<Empty.Title>No resources found</Empty.Title>
						<Empty.Description>
							{searchQuery.trim()
								? 'Try adjusting your search query'
								: 'Get started by adding an application or database'}
						</Empty.Description>
					</Empty.Header>
				</Empty.Root>
			{:else}
				<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
					{#each filteredResources as resource (resource.type === 'application' ? `app-${resource.data.id}` : `db-${resource.data.id}`)}
						{#if resource.type === 'application'}
							<ApplicationCard
								application={resource.data}
								onclick={() => {
									const envId = resource.data.environment_id;
									console.log(envId);
									goto(`/dashboard/project/${projectId}/${envId}/app/${resource.data.id}`);
								}}
							/>
						{:else}
							<DatabaseCard
								database={resource.data}
								onclick={() => {
									const envId = resource.data.environment_id;
									goto(`/dashboard/project/${projectId}/${envId}/db/${resource.data.id}`);
								}}
							/>
						{/if}
					{/each}
				</div>
			{/if}
		</div>
	{/if}
</div>
