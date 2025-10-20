<script lang="ts">
	import '../app.css';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { ModeWatcher } from 'mode-watcher';
	import { browser } from '$app/environment';
	let { children } = $props();

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				enabled: browser
			}
		}
	});
</script>

<ModeWatcher defaultMode={'dark'} />
<div>
	<QueryClientProvider client={queryClient}>
		{@render children()}
	</QueryClientProvider>
</div>
