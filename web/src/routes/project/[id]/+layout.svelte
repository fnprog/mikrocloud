<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import SiteHeader from '$lib/components/ProjectHeader.svelte';
	import { page } from '$app/state';

	const url = page.url.pathname;
	const projectId = page.params.id;

	const data = [
		{
			title: 'Overview',
			url: `/project/${projectId}`
		},
		{
			title: 'Deployments',
			url: `/project/${projectId}/deployments`
		},
		{
			title: 'Logs',
			url: `/project/${projectId}/logs`
		},
		{
			title: 'Analytics',
			url: `/project/${projectId}/analytics`
		},
		{
			title: 'Environment variables',
			url: `/project/${projectId}/env`
		},
		{
			title: 'Processes',
			url: `/project/${projectId}/processes`
		},
		{
			title: 'Domains',
			url: `/project/${projectId}/domains`
		},
		{
			title: 'Networking',
			url: `/project/${projectId}/networking`
		},
		{
			title: 'Disks',
			url: `/project/${projectId}/disks`
		},
		{
			title: 'Web terminal',
			url: `/project/${projectId}/term`
		},
		{
			title: 'Settings',
			url: `/project/${projectId}/settings`
		}
	];

	let { children } = $props();
</script>

<div class="[--header-height:calc(--spacing(14))]">
	<Sidebar.Provider class="flex flex-col">
		<SiteHeader />
		<div class="flex flex-1">
			<Sidebar.Root class="top-(--header-height) h-[calc(100svh-var(--header-height))]!">
				<Sidebar.Content>
					<Sidebar.Menu class="p-4 space-y-1">
						{#each data as item (item.title)}
							<Sidebar.MenuItem>
								<Sidebar.MenuButton tooltipContent={item.title} isActive={item.url === url}>
									{#snippet child({ props })}
										<a href={item.url} {...props}>
											<span>{item.title}</span>
										</a>
									{/snippet}
								</Sidebar.MenuButton>
							</Sidebar.MenuItem>
						{/each}
					</Sidebar.Menu>
				</Sidebar.Content>
			</Sidebar.Root>
			<Sidebar.Inset>
				{@render children()}
			</Sidebar.Inset>
		</div>
	</Sidebar.Provider>
</div>
