<script lang="ts">
	import { page } from '$app/state';
	import { Building2, Users, Settings } from 'lucide-svelte';

	let { children } = $props();

	const menuItems = [
		{
			label: 'General',
			href: '/dashboard/account/organizations/general',
			icon: Building2
		},
		{
			label: 'Members',
			href: '/dashboard/account/organizations/members',
			icon: Users
		},
		{
			label: 'Settings',
			href: '/dashboard/account/organizations/settings',
			icon: Settings
		}
	];
</script>

<div class="flex min-h-screen">
	<aside class="w-64 border-r border-border bg-card/50">
		<div class="p-6">
			<h2 class="text-lg font-semibold mb-6">Organization</h2>
			<nav class="space-y-1">
				{#each menuItems as item}
					<a
						href={item.href}
						class="flex items-center gap-3 px-3 py-2 rounded-lg transition-colors {page.url
							.pathname === item.href
							? 'bg-primary text-primary-foreground'
							: 'hover:bg-muted'}"
					>
						<item.icon class="w-4 h-4" />
						<span>{item.label}</span>
					</a>
				{/each}
			</nav>
		</div>
	</aside>

	<main class="flex-1 p-8">
		<div class="max-w-4xl mx-auto">
			{@render children()}
		</div>
	</main>
</div>
