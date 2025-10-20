<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import {
		SvelteFlow,
		Controls,
		Background,
		MiniMap,
		type Node,
		type Edge,
		BackgroundVariant
	} from '@xyflow/svelte';
	import '@xyflow/svelte/dist/style.css';
	import FocalpointNode from '$lib/components/ui/project-architecture/FocalpointNode.svelte';
	import { createApplicationsFetchQuery } from '$lib/features/applications/queries';
	import { createDatabasesFetchQuery } from '$lib/features/databases/queries';
	import { createEnvironmentsListQuery } from '$lib/features/environments/queries';
	import type { Application } from '$lib/features/applications/types';
	import type { Database } from '$lib/features/databases/types';
	import { Spinner } from '$lib/components/ui/spinner';

	const projectId = $derived(page.params.id!);

	const applicationsQuery = $derived(createApplicationsFetchQuery(projectId, ''));
	const databasesQuery = $derived(createDatabasesFetchQuery(projectId, ''));
	const environmentsQuery = $derived(createEnvironmentsListQuery(projectId));

	const nodeTypes = {
		focalpointNode: FocalpointNode
	};

	const nodes = $derived.by(() => {
		const apps: Application[] = applicationsQuery.data || [];
		const dbs: Database[] = databasesQuery.data || [];
		const envs = environmentsQuery.data || [];

		const envMap = new Map(envs.map((env) => [env.id, env.name]));

		const result: Node[] = [];
		let xPos = 100;
		let yPos = 100;
		const spacing = 400;
		const rowSpacing = 300;

		apps.forEach((app, index) => {
			result.push({
				id: `app-${app.id}`,
				type: 'focalpointNode',
				position: {
					x: xPos + (index % 3) * spacing,
					y: yPos + Math.floor(index / 3) * rowSpacing
				},
				data: {
					id: app.id,
					name: app.name,
					description: app.description,
					status: app.status,
					nodeType: 'application' as const,
					domain: app.custom_domain || app.generated_domain || app.domain,
					environmentId: app.environment_id,
					environmentName: envMap.get(app.environment_id),
					onViewLogs: (id: string) => {
						goto(`/dashboard/project/${projectId}/${app.environment_id}/app/${id}/logs`);
					},
					onOpenTerminal: (id: string) => {
						goto(`/dashboard/project/${projectId}/${app.environment_id}/app/${id}/terminal`);
					},
					onOpenSettings: (id: string) => {
						goto(`/dashboard/project/${projectId}/${app.environment_id}/app/${id}/general`);
					}
				}
			});
		});

		const dbYOffset = yPos + Math.ceil(apps.length / 3) * rowSpacing + 100;
		dbs.forEach((db, index) => {
			result.push({
				id: `db-${db.id}`,
				type: 'focalpointNode',
				position: {
					x: xPos + (index % 3) * spacing,
					y: dbYOffset + Math.floor(index / 3) * rowSpacing
				},
				data: {
					id: db.id,
					name: db.name,
					description: db.description,
					status: db.status,
					nodeType: 'database' as const,
					dbType: db.type,
					connectionString: db.connection_string,
					environmentId: db.environment_id,
					environmentName: envMap.get(db.environment_id),
					onOpenSettings: (id: string) => {
						goto(`/dashboard/project/${projectId}/${db.environment_id}/db/${id}/overview`);
					}
				}
			});
		});

		return result;
	});

	const edges = $derived.by(() => {
		const result: Edge[] = [];
		return result;
	});

	const isLoading = $derived(
		applicationsQuery.isLoading || databasesQuery.isLoading || environmentsQuery.isLoading
	);
	const hasError = $derived(
		applicationsQuery.isError || databasesQuery.isError || environmentsQuery.isError
	);
</script>

<div class="h-[calc(100vh-200px)] w-full rounded-lg border bg-white">
	{#if isLoading}
		<div class="flex h-full items-center justify-center">
			<div class="flex flex-col items-center gap-2">
				<Spinner class="size-8" />
				<p class="text-muted-foreground text-sm">Loading architecture...</p>
			</div>
		</div>
	{:else if hasError}
		<div class="flex h-full items-center justify-center">
			<div class="text-center">
				<p class="text-destructive mb-2">Failed to load architecture data</p>
			</div>
		</div>
	{:else if nodes.length === 0}
		<div class="flex h-full items-center justify-center">
			<div class="text-center">
				<p class="text-muted-foreground mb-2">No resources to display</p>
				<p class="text-muted-foreground text-sm">
					Add applications or databases to see them in the architecture view
				</p>
			</div>
		</div>
	{:else}
		<SvelteFlow {nodes} {edges} {nodeTypes} fitView>
			<Controls />
			<Background variant={BackgroundVariant.Dots} />
			<MiniMap />
		</SvelteFlow>
	{/if}
</div>
