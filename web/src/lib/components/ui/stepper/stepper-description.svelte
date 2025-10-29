<script lang="ts" module>
	export interface StepperDescriptionProps {
		class?: string;
		children: Snippet;
	}
</script>

<script lang="ts">
	import { getContext } from 'svelte';
	import type { Snippet } from 'svelte';
	import type { StepItemContext } from './stepper.svelte';

	let { class: className = '', children }: StepperDescriptionProps = $props();

	const itemCtx = getContext<StepItemContext>('stepItem');
	if (!itemCtx) throw new Error('StepperDescription must be used within StepperItem');
</script>

<div
	data-slot="stepper-description"
	data-state={itemCtx.state}
	class="text-sm text-muted-foreground {className}"
>
	{@render children()}
</div>
