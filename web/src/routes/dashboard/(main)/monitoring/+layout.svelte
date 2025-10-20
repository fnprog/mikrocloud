<script lang="ts">
	import { page } from '$app/state';
	import { cn } from '$lib/utils';
	import {
		LayoutGrid,
		Box,
		Server,
		Globe,
		Activity,
		Clock,
		Gauge,
		ScrollText,
		FileText,
		ListChecks
	} from 'lucide-svelte';

	const navItems = [
		{ href: '/dashboard/monitoring/overview', label: 'Overview', icon: LayoutGrid },
		{ href: '/dashboard/monitoring/containers', label: 'Containers', icon: Box },
		{ href: '/dashboard/monitoring/servers', label: 'Servers', icon: Server },
		{ href: '/dashboard/monitoring/traefik', label: 'Traefik Proxy', icon: Globe },
		{ href: '/dashboard/monitoring/services', label: 'PaaS Services', icon: Activity },
		{ href: '/dashboard/monitoring/deployments', label: 'Deployments', icon: Clock },
		{ href: '/dashboard/monitoring/uptime', label: 'Uptime Checks', icon: Gauge },
		{ href: '/dashboard/monitoring/apm', label: 'APM Stats', icon: Activity },
		{ href: '/dashboard/monitoring/activities', label: 'Activities', icon: ScrollText },
		{ href: '/dashboard/monitoring/logs', label: 'Logs', icon: FileText },
		{ href: '/dashboard/monitoring/queues', label: 'Queue Metrics', icon: ListChecks }
	];

	const currentPath = $derived(page.url.pathname);
</script>

<div class="flex h-full">
	<aside class="w-64 border-r border-border bg-card/50 overflow-y-auto">
		<nav class="p-4 space-y-1">
			{#each navItems as item}
				<a
					href={item.href}
					class={cn(
						'flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors',
						currentPath === item.href
							? 'bg-primary text-primary-foreground'
							: 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
					)}
				>
					<svelte:component this={item.icon} class="h-4 w-4" />
					{item.label}
				</a>
			{/each}
		</nav>
	</aside>

	<main class="flex-1 overflow-y-auto">
		<slot />
	</main>
</div>
