<script lang="ts" module>
	export interface StepperIndicatorProps {
		class?: string;
		children?: Snippet;
	}
</script>

<script lang="ts">
	import { getContext } from 'svelte';
	import type { Snippet } from 'svelte';
	import type { StepperContext, StepItemContext } from './stepper.svelte';

	let { class: className = '', children }: StepperIndicatorProps = $props();

	const stepperCtx = getContext<StepperContext>('stepper');
	const itemCtx = getContext<StepItemContext>('stepItem');

	if (!stepperCtx || !itemCtx) {
		throw new Error('StepperIndicator must be used within StepperItem');
	}

	const indicatorContent = $derived(
		(itemCtx.isLoading && stepperCtx.indicators.loading) ||
			(itemCtx.state === 'completed' && stepperCtx.indicators.completed) ||
			(itemCtx.state === 'active' && stepperCtx.indicators.active) ||
			(itemCtx.state === 'inactive' && stepperCtx.indicators.inactive)
	);
</script>

<div
	data-slot="stepper-indicator"
	data-state={itemCtx.state}
	class="relative flex items-center overflow-hidden justify-center size-6 shrink-0 border-background bg-accent text-accent-foreground rounded-full text-xs data-[state=completed]:bg-primary data-[state=completed]:text-primary-foreground data-[state=active]:bg-primary data-[state=active]:text-primary-foreground {className}"
>
	<div class="absolute">
		{#if indicatorContent}
			{@render indicatorContent()}
		{:else if children}
			{@render children()}
		{/if}
	</div>
</div>
