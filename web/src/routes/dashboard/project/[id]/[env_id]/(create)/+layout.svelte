<script lang="ts">
	let { children } = $props();
	import { NavBarBottom } from '$lib/components/ui/navbar';
</script>

<NavBarBottom tabs={[]} />
<div class="p-6 pt-[104px]">
	{@render children()}
</div>
