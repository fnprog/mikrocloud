<script lang="ts" module>
	export interface StepperTitleProps {
		class?: string;
		children: Snippet;
	}
</script>

<script lang="ts">
	import { getContext } from 'svelte';
	import type { Snippet } from 'svelte';
	import type { StepItemContext } from './stepper.svelte';

	let { class: className = '', children }: StepperTitleProps = $props();

	const itemCtx = getContext<StepItemContext>('stepItem');
	if (!itemCtx) throw new Error('StepperTitle must be used within StepperItem');
</script>

<h3
	data-slot="stepper-title"
	data-state={itemCtx.state}
	class="text-sm font-medium leading-none {className}"
>
	{@render children()}
</h3>
