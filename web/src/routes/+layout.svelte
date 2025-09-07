<script lang="ts">
	import '../app.css';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	let { children } = $props();

	const queryClient = new QueryClient();
</script>

<div>
	<QueryClientProvider client={queryClient}>
		{@render children()}
	</QueryClientProvider>
</div>
