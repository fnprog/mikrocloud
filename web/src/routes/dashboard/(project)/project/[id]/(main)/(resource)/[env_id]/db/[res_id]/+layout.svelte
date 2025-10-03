<script lang="ts">
	let { children } = $props();
	import { NavBarBottom } from '$lib/components/ui/navbar';

	const tabs = [
		{ name: 'Configuration', href: '/dashboard' },
		{ name: 'Logs', href: '/dashboard/projects' },
		{ name: 'Terminal', href: '/dashboard/analytics' },
		{ name: 'DB Studio', href: '/dashboard/settings' }
	];
</script>

<main>
	<NavBarBottom {tabs} />

	<div class="p-12 py-[80px]">
		{@render children()}
	</div>
</main>
