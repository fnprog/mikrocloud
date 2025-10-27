<script lang="ts">
	let { children } = $props();
	import { NavBarBottom } from '$lib/components/ui/navbar';
</script>

<NavBarBottom tabs={[]} />
<div class="pt-[50px]">
	<div class="min-h-screen">
		{@render children()}
	</div>
</div>
