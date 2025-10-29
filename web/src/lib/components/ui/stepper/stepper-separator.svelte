<script lang="ts" module>
	export interface StepperSeparatorProps {
		class?: string;
	}
</script>

<script lang="ts">
	import { getContext } from 'svelte';
	import type { StepItemContext } from './stepper.svelte';

	let { class: className = '' }: StepperSeparatorProps = $props();

	const itemCtx = getContext<StepItemContext>('stepItem');
	if (!itemCtx) throw new Error('StepperSeparator must be used within StepperItem');
</script>

	<div
		data-slot="stepper-separator"
		data-state={itemCtx.state}
		class="m-0.5 rounded-full bg-muted group-data-[orientation=vertical]/stepper-nav:h-12 group-data-[orientation=vertical]/stepper-nav:w-0.5 group-data-[orientation=horizontal]/stepper-nav:h-0.5 group-data-[orientation=horizontal]/stepper-nav:flex-1 group-data-[orientation=horizontal]/stepper-nav:self-start {className}"
	/>
