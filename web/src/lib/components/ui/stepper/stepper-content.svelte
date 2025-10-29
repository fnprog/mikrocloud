<script lang="ts" module>
	export interface StepperContentProps {
		value: number;
		forceMount?: boolean;
		class?: string;
		children: Snippet;
	}
</script>

<script lang="ts">
	import { getContext } from 'svelte';
	import type { Snippet } from 'svelte';
	import type { StepperContext } from './stepper.svelte';

	let {
		value,
		forceMount = false,
		class: className = '',
		children
	}: StepperContentProps = $props();

	const ctx = getContext<StepperContext>('stepper');
	if (!ctx) throw new Error('StepperContent must be used within Stepper');

	const isActive = $derived(value === ctx.activeStep);
</script>

{#if forceMount || isActive}
	<div
		data-slot="stepper-content"
		data-state={ctx.activeStep}
		class="w-full {className} {!isActive && forceMount ? 'hidden' : ''}"
		hidden={!isActive && forceMount}
	>
		{@render children()}
	</div>
{/if}
