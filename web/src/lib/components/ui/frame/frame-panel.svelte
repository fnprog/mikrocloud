<script lang="ts">
	import type { HTMLAttributes } from "svelte/elements";
	import { cn, type WithElementRef } from "$lib/utils.js";

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<div
	bind:this={ref}
	data-slot="frame-panel"
	class={cn(
		"relative bg-clip-padding not-has-[table]:rounded-xl not-has-[table]:border not-has-[table]:bg-card not-has-[table]:p-5 not-has-[table]:shadow-xs before:pointer-events-none before:absolute before:inset-0 before:rounded-[calc(var(--radius-xl)-1px)] before:shadow-[0_1px_--theme(--color-black/4%)] has-[table]:before:hidden dark:bg-clip-border dark:before:shadow-[0_-1px_--theme(--color-white/8%)]",
		className
	)}
	{...restProps}
>
	{@render children?.()}
</div>
