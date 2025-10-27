<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { Settings, Shield, RefreshCw, Network, Mail } from 'lucide-svelte';

	const navItems = [
		{ path: 'general', label: 'General', icon: Settings },
		// { path: 'advanced', label: 'Advanced', icon: Shield },
		{ path: 'updates', label: 'Updates', icon: RefreshCw },
		{ path: 'smtp', label: 'SMTP', icon: Mail },
		// { path: 'backup-location', label: 'S3 Destinations', icon: Database },
		// { path: 'backup', label: 'Backup', icon: Database },
		// { path: 'oauth-keys', label: 'Oauth Keys', icon: Database },
		{ path: 'tunnels', label: 'Tunnels', icon: Network }
	];

	let { children } = $props();

	const isActive = (path: string) => {
		return (
			page.url.pathname.endsWith(`/${path}`) ||
			(page.url.pathname.endsWith('/settings') && path === 'general')
		);
	};
</script>

<div class="min-h-screen">
	<div class="border-b">
		<div class="max-w-7xl mx-auto px-6">
			<h1 class="text-3xl font-semibold my-[40px] w-full mx-auto">Settings</h1>
		</div>
	</div>
	<div class="mt-[46px]"></div>

	<div class="flex max-w-7xl mx-auto px-6">
		<aside class=" flex gap-6">
			<nav class="w-56 shrink-0">
				<div class="space-y-1">
					{#each navItems as item}
						<button
							onclick={() => goto(`/dashboard/settings/${item.path}`)}
							class="w-full flex items-center gap-3 px-4 py-2 text-sm rounded-lg transition-colors {isActive(
								item.path
							)
								? 'bg-accent text-accent-foreground font-medium'
								: 'text-muted-foreground hover:bg-accent/50'}"
						>
							<item.icon class="h-4 w-4" />
							{item.label}
						</button>
					{/each}
				</div>
			</nav>
		</aside>
		<main class="flex-1">
			<div class="max-w-4xl mx-auto">
				{@render children()}
			</div>
		</main>
	</div>
</div>
