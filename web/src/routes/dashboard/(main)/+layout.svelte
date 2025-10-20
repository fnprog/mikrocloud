<script lang="ts">
	let { children } = $props();
	import { NavBarBottom } from '$lib/components/ui/navbar';

	const tabs = [
		{ name: 'Overview', href: '/dashboard/overview' },
		{ name: 'Projects', href: '/dashboard/projects' },
		{ name: 'Infrastructure', href: '/dashboard/infra' },
		// { name: 'Observability', href: '/dashboard/monitoring' },
		{ name: 'Version Control', href: '/dashboard/git' },
		// { name: 'Notifications', href: '/dashboard/notifications/email' },
		// { name: 'Organization', href: '/dashboard/organizations/general' },
		// { name: 'Agent', href: '/dashboard/agent' },
		{ name: 'Settings', href: '/dashboard/settings' }
	];
</script>

<main>
	<NavBarBottom {tabs} />

	<div class="p-6 pt-[104px]">
		{@render children()}
	</div>
</main>
