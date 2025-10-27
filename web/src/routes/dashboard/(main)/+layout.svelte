<script lang="ts">
	let { children } = $props();
	import { NavBarBottom } from '$lib/components/ui/navbar';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	const tabs = [
		{ name: 'Overview', href: '/dashboard/overview' },
		{ name: 'Projects', href: '/dashboard/projects' },
		// { name: 'Infrastructure', href: '/dashboard/infra' },
		{ name: 'Version Control', href: '/dashboard/git' },
		// { name: 'Organization', href: '/dashboard/organizations/general' },
		{ name: 'Settings', href: '/dashboard/settings' }
	];

	onMount(async () => {
		try {
			const response = await fetch('/api/auth/setup');
			const data = await response.json();
			if (!data.is_setup) {
				goto('/register?firstTime=true');
			}
		} catch (error) {
			console.error('Failed to check setup status:', error);
		}
	});
</script>

<main>
	<NavBarBottom {tabs} />

	<div class="pt-[50px]">
		{@render children()}
	</div>
</main>
