<script lang="ts">
	import { page } from '$app/stores';
	import { createQuery } from '@tanstack/svelte-query';
	import { projectsApi, environmentsApi, databasesApi } from '$lib/api';
	import { Loader2 } from 'lucide-svelte';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb';

	const projectId = $derived($page.params.id);
	const envId = $derived($page.params.env_id);
	const resId = $derived($page.params.res_id);
	const resourceType = $derived($page.url.pathname.includes('/db/') ? 'database' : 'application');

	const projectQuery = createQuery({
		queryKey: ['project', projectId],
		queryFn: () => projectsApi.get(projectId!),
		enabled: !!projectId
	});

	const environmentQuery = createQuery({
		queryKey: ['environment', projectId, envId],
		queryFn: () => environmentsApi.get(projectId!, envId!),
		enabled: !!projectId && !!envId
	});

	const databaseQuery = createQuery({
		queryKey: ['database', projectId, resId],
		queryFn: () => databasesApi.get(projectId!, resId!),
		enabled: !!projectId && !!resId && resourceType === 'database'
	});

	const breadcrumbs = $derived.by(() => {
		const crumbs = [];

		if (projectId) {
			crumbs.push({
				label: 'My Projects',
				href: '/dashboard'
			});
		}

		if (projectId && $projectQuery.data) {
			crumbs.push({
				label: $projectQuery.data.name,
				href: `/dashboard/project/${projectId}`
			});
		} else if (projectId) {
			crumbs.push({
				label: '...',
				href: null,
				loading: true
			});
		}

		if (envId && $environmentQuery.data) {
			crumbs.push({
				label: $environmentQuery.data.name,
				href: `/dashboard/project/${projectId}/${envId}`
			});
		} else if (envId) {
			crumbs.push({
				label: '...',
				href: null,
				loading: true
			});
		}

		if (resId && resourceType === 'database' && $databaseQuery.data) {
			crumbs.push({
				label: $databaseQuery.data.name,
				href: `/dashboard/project/${projectId}/${envId}/db/${resId}`
			});
		} else if (resId && resourceType === 'database') {
			crumbs.push({
				label: '...',
				href: null,
				loading: true
			});
		}

		return crumbs;
	});
</script>

{#if breadcrumbs.length > 0}
	<Breadcrumb.Root>
		<Breadcrumb.List>
			{#each breadcrumbs as crumb, i (i)}
				<Breadcrumb.Item>
					{#if crumb.href}
						<Breadcrumb.Link href={crumb.href}>
							{crumb.label}
						</Breadcrumb.Link>
					{:else}
						<Breadcrumb.Page class="flex items-center gap-1.5">
							{#if crumb.loading}
								<Loader2 class="w-3 h-3 animate-spin" />
							{/if}
							{crumb.label}
						</Breadcrumb.Page>
					{/if}
				</Breadcrumb.Item>
				{#if i < breadcrumbs.length - 1}
					<Breadcrumb.Separator>/</Breadcrumb.Separator>
				{/if}
			{/each}
		</Breadcrumb.List>
	</Breadcrumb.Root>
{/if}
