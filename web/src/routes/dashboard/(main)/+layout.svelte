<script lang="ts">
	let { children } = $props();
	import { NavBarBottom } from '$lib/components/ui/navbar';

	const tabs = [
		{ name: 'Overview', href: '/dashboard' },
		{ name: 'Projects', href: '/dashboard/projects' },
		{ name: 'Analytics', href: '/dashboard/analytics' },
		{ name: 'Settings', href: '/dashboard/settings' }
	];
</script>

<main>
	<NavBarBottom {tabs} />

	<div class="p-6 pt-[104px]">
		{@render children()}
	</div>
</main>
