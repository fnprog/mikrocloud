<script lang="ts">
	let { children } = $props();
	import { page } from '$app/state';
	import { NavBarBottom } from '$lib/components/ui/navbar';

	const projectId = $derived(page.params.id);

	const tabs = $derived([
		{ name: 'Overview', href: `/dashboard/project/${projectId}/overview` },
		// { name: 'Architecture', href: `/dashboard/project/${projectId}/schema` },
		// { name: 'Shared Variables', href: `/dashboard/project/${projectId}/shared-vars` },
		{ name: 'Settings', href: `/dashboard/project/${projectId}/settings` }
	]);
</script>

<NavBarBottom {tabs} />
<div class="pt-[50px]">
	<div class="min-h-screen">
		{@render children()}
	</div>
</div>
