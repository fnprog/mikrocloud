<script lang="ts">
	import { page } from '$app/state';
	import { LoaderCircle } from 'lucide-svelte';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb';
	import OrgSwitcher from './org-switcher.svelte';
	import ProjectSwitcher from './project-switcher.svelte';
	import ResourceSwitcher from './resource-switcher.svelte';
	import { createOrganizationsListQuery } from '$lib/features/organizations/queries';
	import { createProjectQuery } from '$lib/features/projects/queries';
	import { createEnvironmentQuery } from '$lib/features/environments/queries';
	import { createDatabaseFetchQuery } from '$lib/features/databases/queries';
	import { createApplicationFetchQuery } from '$lib/features/applications/queries';

	const projectId = $derived(page.params.id);
	const envId = $derived(page.params.env_id);
	const resId = $derived(page.params.res_id);

	const resourceType = $derived(
		page.url.pathname.includes('/db/')
			? 'database'
			: page.url.pathname.includes('/app/')
				? 'application'
				: null
	);

	const orgQuery = createOrganizationsListQuery();
	const projectQuery = $derived(createProjectQuery(projectId));
	const environmentQuery = $derived(createEnvironmentQuery(projectId, envId));
	const databaseQuery = $derived(
		createDatabaseFetchQuery(projectId, resourceType === 'database' ? resId : '')
	);
	const applicationQuery = $derived(
		createApplicationFetchQuery(projectId, envId, resourceType === 'application' ? resId : '')
	);

	const currentOrg = $derived(orgQuery.data?.[0]);
	const shouldShowBreadcrumbs = $derived(true);
</script>

{#if shouldShowBreadcrumbs}
	<Breadcrumb.Root>
		<Breadcrumb.List>
			<Breadcrumb.Item>
				{#if currentOrg}
					<OrgSwitcher currentOrgId={currentOrg.id}>
						{currentOrg.name}
					</OrgSwitcher>
				{:else}
					<span class="flex items-center gap-1.5">
						<LoaderCircle class="w-3 h-3 animate-spin" />
						Loading...
					</span>
				{/if}
			</Breadcrumb.Item>

			{#if projectId}
				<Breadcrumb.Separator>/</Breadcrumb.Separator>
				<Breadcrumb.Item>
					{#if projectQuery.data}
						<ProjectSwitcher currentProjectId={projectId}>
							{projectQuery.data.name}
						</ProjectSwitcher>
					{:else}
						<span class="flex items-center gap-1.5">
							<LoaderCircle class="w-3 h-3 animate-spin" />
							...
						</span>
					{/if}
				</Breadcrumb.Item>
			{/if}

			{#if envId && environmentQuery.data}
				<Breadcrumb.Separator>/</Breadcrumb.Separator>
				<Breadcrumb.Item>
					<Breadcrumb.Page>{environmentQuery.data.name}</Breadcrumb.Page>
				</Breadcrumb.Item>
			{/if}

			{#if resId && resourceType && projectId}
				<Breadcrumb.Separator>/</Breadcrumb.Separator>
				<Breadcrumb.Item>
					{#if resourceType === 'database' && databaseQuery.data}
						<ResourceSwitcher
							environmentId={envId}
							{projectId}
							currentResourceId={resId}
							currentResourceType="database"
						>
							{databaseQuery.data.name}
						</ResourceSwitcher>
					{:else if resourceType === 'application' && applicationQuery.data}
						<ResourceSwitcher
							{projectId}
							environmentId={envId}
							currentResourceId={resId}
							currentResourceType="application"
						>
							{applicationQuery.data.name}
						</ResourceSwitcher>
					{:else}
						<span class="flex items-center gap-1.5">
							<LoaderCircle class="w-3 h-3 animate-spin" />
							...
						</span>
					{/if}
				</Breadcrumb.Item>
			{/if}
		</Breadcrumb.List>
	</Breadcrumb.Root>
{/if}
