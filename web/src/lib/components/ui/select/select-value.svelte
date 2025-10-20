<script lang="ts">
	import { getSelectValueContext } from "bits-ui";
	import { cn } from "$lib/utils.js";

	let {
		class: className,
		placeholder = "Select an option",
		...restProps
	}: {
		class?: string;
		placeholder?: string;
		[key: string]: any;
	} = $props();

	const ctx = getSelectValueContext();
	const selectedLabel = $derived(ctx.selectedLabel || placeholder);
</script>

<span
	data-slot="select-value"
	class={cn("pointer-events-none", className)}
	{...restProps}
>
	{selectedLabel}
</span>
