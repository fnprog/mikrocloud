<script lang="ts">
	let { children } = $props();
	import { NavBarTop } from '$lib/components/ui/navbar';
</script>

<div class="min-h-screen">
	<NavBarTop />

	{@render children()}
</div>
