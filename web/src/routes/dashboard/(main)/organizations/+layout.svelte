<script lang="ts">
	import { page } from '$app/state';
	import { Building2, Users, Settings } from 'lucide-svelte';

	let { children } = $props();

	const menuItems = [
		{
			label: 'General',
			href: 'general',
			icon: Building2
		},
		{
			label: 'Members',
			href: 'members',
			icon: Users
		},
		{
			label: 'Settings',
			href: 'settings',
			icon: Settings
		}
	];

	const isActive = (path: string) => {
		return (
			page.url.pathname.endsWith(`/${path}`) ||
			(page.url.pathname.endsWith('/organizations') && path === 'general')
		);
	};
</script>

<div class="container max-w-7xl py-8">
	<aside class="flex gap-6">
		<nav class="w-56 flex-shrink-0">
			<div class="space-y-1">
				{#each menuItems as item}
					<a
						href={`/dashboard/organizations/${item.href}`}
						class="w-full flex items-center gap-3 px-4 py-2 text-sm rounded-lg transition-colors {isActive(
							item.href
						)
							? 'bg-accent text-accent-foreground font-medium'
							: 'text-muted-foreground hover:bg-accent/50'}"
					>
						<item.icon class="h-4 w-4" />
						{item.label}
					</a>
				{/each}
			</div>
		</nav>

		<div class="flex-1">
			{@render children()}
		</div>
	</aside>
</div>
